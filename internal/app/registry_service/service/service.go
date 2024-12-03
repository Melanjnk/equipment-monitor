package service

import "github.com/Melanjnk/equipment-monitor/internal/app/registry_service/dtos"

type EquipmentRepository interface {
	CreateOne(equipmentCreate *dtos.EquipmentCreate) (string, error)
	CreateMany(equipmentCreate []dtos.EquipmentCreate) ([]string, error)
	UpdateById(equipmentUpdate *dtos.EquipmentUpdate, id string) (bool, error)
	UpdateByIds(equipmentUpdate *dtos.EquipmentUpdate, ids []string) ([]string, error)
	UpdateByConditions(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) ([]string, error)
	DeleteById(id string) (bool, error)
	DeleteByIds(ids []string) ([]string, error)
	DeleteByConditions(equipmentFilter *dtos.EquipmentFilter) ([]string, error)
	FindById(id string) (*dtos.EquipmentGet, error)
	FindByIds(ids []string) ([]dtos.EquipmentGet, error)
	FindByConditions(equipmentFilter *dtos.EquipmentFilter) ([]dtos.EquipmentGet, error)
}

type Equipment struct {
	repository EquipmentRepository
}

func NewEquipment(repository EquipmentRepository) Equipment {
	return Equipment{repository: repository}
}

func (service Equipment) CreateOne(equipmentCreate *dtos.EquipmentCreate) (string, error) {
	return service.repository.CreateOne(equipmentCreate)
}

func (service Equipment) CreateMany(equipmentCreate []dtos.EquipmentCreate) ([]string, error) {
	return service.repository.CreateMany(equipmentCreate)
}

func (service Equipment) UpdateById(equipmentUpdate *dtos.EquipmentUpdate, id string) (bool, error) {
	return service.repository.UpdateById(equipmentUpdate, id)
}

func (service Equipment) UpdateByIds(equipmentUpdate *dtos.EquipmentUpdate, ids []string) ([]string, error) {
	return service.repository.UpdateByIds(equipmentUpdate, ids)
}

func (service Equipment) UpdateByConditions(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	return service.repository.UpdateByConditions(equipmentUpdate, equipmentFilter)
}

func (service Equipment) DeleteById(id string) (bool, error) {
	return service.repository.DeleteById(id)
}

func (service Equipment) DeleteByIds(ids []string) ([]string, error) {
	return service.repository.DeleteByIds(ids)
}

func (service Equipment) DeleteByConditions(equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	return service.repository.DeleteByConditions(equipmentFilter)
}

func (service Equipment) FindById(id string) (*dtos.EquipmentGet, error) {
	return service.repository.FindById(id)
}

func (service Equipment) FindByIds(ids []string) ([]dtos.EquipmentGet, error) {
	return service.repository.FindByIds(ids)
}

func (service Equipment) FindByConditions(equipmentFilter *dtos.EquipmentFilter) ([]dtos.EquipmentGet, error) {
	return service.repository.FindByConditions(equipmentFilter)
}
