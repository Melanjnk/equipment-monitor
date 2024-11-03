package service

import (
	"github.com/gofrs/uuid"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/dtos"
)

type EquipmentRepo interface {
	List() ([]*dtos.EquipmentGet, error)
	Create(eqc *dtos.EquipmentCreate) (uuid.UUID, error)
	Update(equ *dtos.EquipmentUpdate) (bool, error)
	FindById(id uuid.UUID) (*dtos.EquipmentGet, error)
	RemoveById(id uuid.UUID) (bool, error)
}

type Equipment struct {
	repo EquipmentRepo
}

func NewEquipment(repo EquipmentRepo) Equipment {
	return Equipment{repo: repo}
}

func (service *Equipment) List() ([]*dtos.EquipmentGet, error) {
	return service.repo.List()
}

func (service *Equipment) Create(eqc *dtos.EquipmentCreate) (uuid.UUID, error) {
	return service.repo.Create(eqc)
}

func (service *Equipment) Update(equ *dtos.EquipmentUpdate) (bool, error) {
	return service.repo.Update(equ)
}

func (service *Equipment) Get(ids string) (*dtos.EquipmentGet, error) {
	id, err := uuid.FromString(ids)
	if err != nil {
		return nil, err
	}

	return service.repo.FindById(id)
}

func (service *Equipment) Delete(eqId string) (bool, error) {
	id, err := uuid.FromString(eqId)
	if err != nil {
		return false, err
	}
	return service.repo.RemoveById(id)
}
