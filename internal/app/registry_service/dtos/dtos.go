package dtos

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/schema"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/model"
)

const (
	emptyParameters string = "Equipment parameters should not be empty"
	invalidFieldValue = "Invalid equipment %s value: %d"
	cannotBeUsedTogether = "Fields `%s` and `%s` cannot be used together"
	mustPrecede = "%s must precede %s"
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
	if equipmentUpdate.Status != nil && !equipmentUpdate.Status.IsValid() {
		return fmt.Errorf("Invalid Equipment.Status: %d", *equipmentUpdate.Status)
	}
	// TODO: check non-empty Parameters
	return nil
}


type EquipmentGet = model.Equipment


type EquipmentFilter struct {
	Kinds			[]model.EquipmentKind		`schema:"kind"`
	NoKinds			[]model.EquipmentKind		`schema:"no_kind"`
	Statuses		[]model.OperationalStatus	`schema:"status"`
	NoStatuses		[]model.OperationalStatus	`schema:"no_status"`
	CreatedSince	*time.Time					`schema:"created_since"`
	CreatedUntil	*time.Time					`schema:"created_until"`
	UpdatedSince	*time.Time					`schema:"updated_since"`
	UpdatedUntil	*time.Time					`schema:"updated_until"`
}

func (equipmentFilter *EquipmentFilter) Validate() error {
	if equipmentFilter.Kinds != nil {
		if equipmentFilter.NoKinds != nil {
			return fmt.Errorf(cannotBeUsedTogether, "kind", "no_kind")
		}
		for _, kind := range equipmentFilter.Kinds {
			if !kind.IsValid() {
				return fmt.Errorf(invalidFieldValue, "kind", kind)
			}
		}
	} else if equipmentFilter.Kinds != nil {
		for _, kind := range equipmentFilter.NoKinds {
			if !kind.IsValid() {
				return fmt.Errorf(invalidFieldValue, "no_kind", kind)
			}
		}
	}
	
	if equipmentFilter.Statuses != nil {
		if equipmentFilter.NoStatuses != nil {
			return fmt.Errorf(cannotBeUsedTogether, "status", "no_status")
		}
		for _, status := range equipmentFilter.Statuses {
			if !status.IsValid() {
				return fmt.Errorf(invalidFieldValue, "status", status)
			}
		}
	} else if equipmentFilter.NoStatuses != nil {
		for _, status := range equipmentFilter.NoStatuses {
			if !status.IsValid() {
				return fmt.Errorf(invalidFieldValue, "no_status", status)
			}
		}
	}

	if equipmentFilter.CreatedSince != nil {
		if equipmentFilter.CreatedUntil != nil && equipmentFilter.CreatedSince.After(*equipmentFilter.CreatedUntil) {
			return fmt.Errorf(mustPrecede, "`created_since`", "`created_until`")
		}
		if equipmentFilter.UpdatedSince != nil && equipmentFilter.UpdatedSince != nil && equipmentFilter.UpdatedSince.Before(*equipmentFilter.CreatedSince) {
			return fmt.Errorf(mustPrecede, "`created_since`", "`updated_since`")
		}
	}
	if equipmentFilter.UpdatedSince != nil && equipmentFilter.UpdatedUntil != nil && equipmentFilter.UpdatedUntil.Before(*equipmentFilter.UpdatedSince) {
		return fmt.Errorf(mustPrecede, "`updated_since`", "`updated_until`")
	}
	return nil
}

func EquipmentFilterFromRequest(request *http.Request) (*EquipmentFilter, error) {
	var err error
	if err = request.ParseForm(); err == nil {
		var equipmentFilter EquipmentFilter
		if err = schema.NewDecoder().Decode(&equipmentFilter, request.Form); err == nil {
			if err = equipmentFilter.Validate(); err == nil {
				return &equipmentFilter, nil
			}
		}
	}
	return nil, err
}
