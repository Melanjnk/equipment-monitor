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
	eitherAtOrSinceBefore = "Either `%[1]s_at` or `%[1]s_since`/`%[1]s_before` must me used"
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
	Statuses		[]model.OperationalStatus	`schema:"status"`
	CreatedAt		*time.Time					`schema:"created_at"`
	CreatedSince	*time.Time					`schema:"created_since"`
	CreatedBefore	*time.Time					`schema:"created_before"`
	UpdatedAt		*time.Time					`schema:"updated_at"`
	UpdatedSince	*time.Time					`schema:"updated_since"`
	UpdatedBefore	*time.Time					`schema:"updated_before"`
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
		for _, kind := range equipmentFilter.Kinds {
			if !kind.IsValid() {
				return fmt.Errorf(invalidFieldValue, "kind", kind)
			}
		}
	}
	if equipmentFilter.Statuses != nil {
		for _, status := range equipmentFilter.Kinds {
			if !status.IsValid() {
				return fmt.Errorf(invalidFieldValue, "status", status)
			}
		}
	}
	if equipmentFilter.CreatedAt != nil {
		if equipmentFilter.CreatedSince != nil || equipmentFilter.CreatedBefore != nil {
			return fmt.Errorf(eitherAtOrSinceBefore, "created")
		}
		if !precedesOthers(*equipmentFilter.CreatedAt,
			equipmentFilter.UpdatedAt, equipmentFilter.UpdatedSince, equipmentFilter.UpdatedBefore,
		) {
			return fmt.Errorf(mustPrecede, "`created_at`", "`updated_at`/`updated_since`/`updated_before`")
		}
	} else if equipmentFilter.CreatedSince != nil {
		if !precedesOthers(*equipmentFilter.CreatedSince,
			equipmentFilter.CreatedBefore, equipmentFilter.UpdatedAt, equipmentFilter.UpdatedSince, equipmentFilter.UpdatedBefore,
		) {
			return fmt.Errorf(mustPrecede, "`created_since`", "`created_before`/`updated_at`/`updated_since`/`updated_before`")
		}
	} else if equipmentFilter.CreatedBefore != nil {
		if !precedesOthers(*equipmentFilter.CreatedBefore,
			equipmentFilter.UpdatedAt, equipmentFilter.UpdatedSince, equipmentFilter.UpdatedBefore,
		) {
			return fmt.Errorf(mustPrecede, "`created_before`", "`updated_at`/`updated_since`/`updated_before`")
		}
	}
	if equipmentFilter.UpdatedAt != nil {
		if equipmentFilter.UpdatedSince != nil || equipmentFilter.UpdatedBefore != nil {
			return fmt.Errorf(eitherAtOrSinceBefore, "updated")
		}
	} else if equipmentFilter.UpdatedSince != nil {
		if !precedesOthers(*equipmentFilter.UpdatedSince, equipmentFilter.UpdatedBefore) {
			return fmt.Errorf(mustPrecede, "`updated_since`", "`updated_before`")
		}
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
