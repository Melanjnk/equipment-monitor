package service

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/dtos"
)

const failedToParseUUID string = "%s: failed to parse `%s` as UUID: %v"

type EquipmentRepository interface {
	List() ([]*dtos.EquipmentGet, error)
	Create(equipmentCreate *dtos.EquipmentCreate) (uuid.UUID, error)
	Update(equipmentUpdate *dtos.EquipmentUpdate) (bool, error)
	FindById(id uuid.UUID) (*dtos.EquipmentGet, error)
	RemoveById(id uuid.UUID) (bool, error)
}

type Equipment struct {
	repository EquipmentRepository
}

func NewEquipment(repository EquipmentRepository) Equipment {
	return Equipment{repository: repository}
}

func (service *Equipment) List() ([]*dtos.EquipmentGet, error) {
	return service.repository.List()
}

func (service *Equipment) Create(equipmentCreate *dtos.EquipmentCreate) (uuid.UUID, error) {
	return service.repository.Create(equipmentCreate)
}

func (service *Equipment) Update(equipmentUpdate *dtos.EquipmentUpdate) (bool, error) {
	return service.repository.Update(equipmentUpdate)
}

func (service *Equipment) Get(equipmentId string) (*dtos.EquipmentGet, error) {
	id, err := uuid.FromString(equipmentId)
	if err != nil {
		return nil, fmt.Errorf(failedToParseUUID, "Get", equipmentId, err)
	}
	return service.repository.FindById(id)
}

func (service *Equipment) Delete(equipmentId string) (bool, error) {
	id, err := uuid.FromString(equipmentId)
	if err != nil {
		return false, fmt.Errorf(failedToParseUUID, "Delete", equipmentId, err)
	}
	return service.repository.RemoveById(id)
}
