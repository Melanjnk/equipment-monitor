package dtos

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"github.com/gofrs/uuid"
	"github.com/gorilla/schema"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/model"
)

const (
	emptyParameters string = "Equipment parameters should not be empty"
	invalidFieldValue = "Invalid equipment %s value: %d"
	/*eitherAtOrSinceAfterBeforeUntil = "Either `%[1]s_at` or `%[1]s_since`/`%[1]s_after`/`%[1]s_before/`%[1]s_until` can be used"
	eitherSinceOrAfter = "Either `%[1]s_since` or `%[1]s_after` can be used"
	eitherBeforeOrUntil = "Either `%[1]s_before` or `%[1]s_until` can be used"*/
	cannotBeUsedTogether = "Fields `%s` and `%s` cannot be used together"
	mustPrecede = "%s must precede %s"
)

type validableDTO interface {
	Validate() error
}

func FromJSON[DTO validableDTO](decoder *json.Decoder) (*DTO, error) {
	var dto DTO
	err := decoder.Decode(&dto)
	if err == nil {
		err = dto.Validate()
		if err == nil {
			return &dto, nil
		}
	}
	return nil, err
}

func FromRequestJSON[DTO validableDTO](request *http.Request) (*DTO, error) {
	return FromJSON[DTO](json.NewDecoder(request.Body))
}


type EquipmentCreate struct {
	Kind		model.EquipmentKind		`json:"kind"`
	Parameters	map[string]interface{}	`json:"parameters"`
}

func (eqc EquipmentCreate) Validate() error {
	if !eqc.Kind.IsValid() {
		return fmt.Errorf("Invalid Equipment.Kind: %d", eqc.Kind)
	}
	if eqc.Parameters == nil { // TODO: check when not nil
		return errors.New(emptyParameters)
	}
	return nil
}


type EquipmentUpdate struct {
	Id			uuid.UUID					`json:"id"`
	Status		*model.OperationalStatus	`json:"status"`
	Parameters	*map[string]interface{}		`json:"parameters"`
}

func (equipmentUpdate EquipmentUpdate) Validate() error {
	if equipmentUpdate.Status != nil && !equipmentUpdate.Status.IsValid() {
		return fmt.Errorf("Invalid Equipment.Status: %d", *equipmentUpdate.Status)
	}
	// TODO: check non-empty Parameters
	return nil
}


type EquipmentGet struct {
	Id			uuid.UUID				`json:"id"`
	Kind		model.EquipmentKind		`json:"kind"`
	Status		model.OperationalStatus	`json:"status"`
	Parameters	map[string]interface{}	`json:"parameters"`
	CreatedAt	time.Time				`json:"created_at"`
	UpdatedAt	time.Time				`json:"updated_at"`
}

func EquipmentGetFromModel(equipmentModel model.Equipment) *EquipmentGet {
	var parameters map[string]interface{}
	_ = json.Unmarshal(equipmentModel.Parameters, &parameters)
	return &EquipmentGet {
		Id:			equipmentModel.Id,
		Kind:		equipmentModel.Kind,
		Status:		equipmentModel.Status,
		Parameters:	parameters,
		CreatedAt:	equipmentModel.CreatedAt,
		UpdatedAt:	equipmentModel.UpdatedAt,
	}
}


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

// precedesOthers returns true if time0 is before or equal to any other non-nil time from params;
// otherwise returns false.
func precedesOthers(time0 time.Time, times ...*time.Time) bool {
	for _, time := range times {
		if time != nil && time0.After(*time) {
			return false
		}
	}
	return true
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
