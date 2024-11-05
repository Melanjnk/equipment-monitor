package dtos

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"github.com/gofrs/uuid"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/model"
)

const emptyParameters string = "Equipment parameters should not be empty"

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
