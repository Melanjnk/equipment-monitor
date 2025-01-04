package dtos

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
	"github.com/gorilla/schema"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/model"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/stringset"
)

const (
	emptyParameters = "Equipment parameters should not be empty"
	invalidFieldValue = "Invalid equipment %s value: %d"
	cannotBeUsedTogether = "Fields `%s` and `%s` cannot be used together"
	mustPrecede = "`%s` must precede `%s`"
)


type EquipmentCreate struct {
	Kind		model.EquipmentKind			`json:"kind" db:"kind,not null"`
	Parameters	model.EquipmentParameters	`json:"parameters" db:"parameters,not null,type:jsonb"`
}

func (equipmentCreate EquipmentCreate) Validate() error {
	if !equipmentCreate.Kind.IsValid() {
		return fmt.Errorf("Invalid Equipment.Kind: %d", equipmentCreate.Kind)
	}
	if equipmentCreate.Parameters == nil { // TODO: check when not nil
		return errors.New(emptyParameters)
	}
	return nil
}


type EquipmentUpdate struct {
	Status		*model.OperationalStatus	`json:"status"`
	Parameters	*model.EquipmentParameters	`json:"parameters" db:"parameters,type:jsonb"`
}

func (equipmentUpdate EquipmentUpdate) Validate() error {
	if equipmentUpdate.Status == nil {
		if equipmentUpdate.Parameters == nil {
			return fmt.Errorf("No parameters to update.")
		}
	} else if !equipmentUpdate.Status.IsValid() {
		return fmt.Errorf("Invalid Equipment.Status: %d", *equipmentUpdate.Status)
	}
	// TODO: check non-empty Parameters
	return nil
}


type EquipmentGet = model.Equipment


type EquipmentFilter struct {
	Ids				[]string					`schema:"id"`
	NoIds			[]string					`schema:"~id"`
	Kinds			[]model.EquipmentKind		`schema:"kind"`
	NoKinds			[]model.EquipmentKind		`schema:"~kind"`
	Statuses		[]model.OperationalStatus	`schema:"status"`
	NoStatuses		[]model.OperationalStatus	`schema:"~status"`
	CreatedSince	*time.Time					`schema:"created_since"`
	CreatedUntil	*time.Time					`schema:"created_until"`
	UpdatedSince	*time.Time					`schema:"updated_since"`
	UpdatedUntil	*time.Time					`schema:"updated_until"`
	Sort			[]string					`schema:"sort"`
	SortMask		uint8						`schema:"-"`
	Limit			*uint						`schema:"limit"`
	Offset			*uint						`schema:"offset"`
}

func IsValidUUID(str string) bool {
	characters := []byte(str)
	i := len(characters)
	if i != 36 {
		return false
	}
	for i > 24 {
		i--
		if !isHexadecimalDigit(characters[i]) {
			return false
		}
	}
	if i--; characters[i] != '-' {
		return false
	}
	for i > 19 {
		i--
		if !isHexadecimalDigit(characters[i]) {
			return false
		}
	}
	if i--; characters[i] != '-' {
		return false
	}
	for i > 14 {
		i--
		if !isHexadecimalDigit(characters[i]) {
			return false
		}
	}
	if i--; characters[i] != '-' {
		return false
	}
	for i > 9 {
		i--
		if !isHexadecimalDigit(characters[i]) {
			return false
		}
	}
	if i--; characters[i] != '-' {
		return false
	}
	for i > 0 {
		i--
		if !isHexadecimalDigit(characters[i]) {
			return false
		}
	}
	return true
}

func (equipmentFilter *EquipmentFilter) Validate() error {
	if equipmentFilter.Ids != nil {
		if equipmentFilter.NoIds != nil {
			return fmt.Errorf(cannotBeUsedTogether, "id", "~id")
		}
		for _, id := range equipmentFilter.Ids {
			if !IsValidUUID(id) {
				return fmt.Errorf(invalidFieldValue, "id", id)
			}
		}
	} else if equipmentFilter.NoIds != nil {
		for _, id := range equipmentFilter.NoIds {
			if !IsValidUUID(id) {
				return fmt.Errorf(invalidFieldValue, "~id", id)
			}
		}
	}

	if equipmentFilter.Kinds != nil {
		if equipmentFilter.NoKinds != nil {
			return fmt.Errorf(cannotBeUsedTogether, "kind", "~kind")
		}
		for _, kind := range equipmentFilter.Kinds {
			if !kind.IsValid() {
				return fmt.Errorf(invalidFieldValue, "kind", kind)
			}
		}
	} else if equipmentFilter.NoKinds != nil {
		for _, kind := range equipmentFilter.NoKinds {
			if !kind.IsValid() {
				return fmt.Errorf(invalidFieldValue, "~kind", kind)
			}
		}
	}
	
	if equipmentFilter.Statuses != nil {
		if equipmentFilter.NoStatuses != nil {
			return fmt.Errorf(cannotBeUsedTogether, "status", "~status")
		}
		for _, status := range equipmentFilter.Statuses {
			if !status.IsValid() {
				return fmt.Errorf(invalidFieldValue, "status", status)
			}
		}
	} else if equipmentFilter.NoStatuses != nil {
		for _, status := range equipmentFilter.NoStatuses {
			if !status.IsValid() {
				return fmt.Errorf(invalidFieldValue, "~status", status)
			}
		}
	}

	if equipmentFilter.CreatedSince != nil {
		if equipmentFilter.CreatedUntil != nil && equipmentFilter.CreatedSince.After(*equipmentFilter.CreatedUntil) {
			return fmt.Errorf(mustPrecede, "created_since", "created_until")
		}
		if equipmentFilter.UpdatedSince != nil && equipmentFilter.UpdatedSince != nil && equipmentFilter.UpdatedSince.Before(*equipmentFilter.CreatedSince) {
			return fmt.Errorf(mustPrecede, "created_since", "updated_since")
		}
	}
	if equipmentFilter.UpdatedSince != nil && equipmentFilter.UpdatedUntil != nil && equipmentFilter.UpdatedUntil.Before(*equipmentFilter.UpdatedSince) {
		return fmt.Errorf(mustPrecede, "updated_since", "updated_until")
	}

	if equipmentFilter.Sort != nil {
		fieldNames := stringset.New(`id`, `kind`, `status`, `created_at`, `updated_at`)
		for i, fieldName := range equipmentFilter.Sort {
			if normalizeFieldName(&fieldName) {
				equipmentFilter.SortMask |= 1 << i
			}
			if len(fieldName) == 0 {
				return fmt.Errorf("Sorting field name #%d is empty", i + 1)
			} else if !fieldNames.Contains(fieldName) {
				return fmt.Errorf("Invalid or duplicate sorting field name: `%s`", fieldName)
			}
			fieldNames.Exclude(fieldName)
			equipmentFilter.Sort[i] = fieldName
		}
	}

	return nil
}

func EquipmentFilterFromIds(ids []string) *EquipmentFilter {
	return &EquipmentFilter{Ids: ids}
}

func EquipmentFilterFromRequest(request *http.Request) (*EquipmentFilter, error) {
	var err error
	if err = request.ParseForm(); err == nil {
		var equipmentFilter EquipmentFilter
		decoder := schema.NewDecoder()
		// Provide splitting paratemer's comma separated value:
		decoder.RegisterConverter([]string{}, func(values string) reflect.Value {
			return reflect.ValueOf(strings.Split(values, `,`))
		})
		if err = decoder.Decode(&equipmentFilter, request.Form); err == nil {
			if err = equipmentFilter.Validate(); err == nil {
				return &equipmentFilter, nil
			}
		}
	}
	return nil, err
}
