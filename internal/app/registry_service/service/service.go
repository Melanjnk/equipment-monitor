package service

import "github.com/Melanjnk/equipment-monitor/internal/app/registry_service/dtos"

type EquipmentRepository interface {
	CreateOne(*dtos.EquipmentCreate) (string, error)
	CreateMany([]dtos.EquipmentCreate) ([]string, error)
	UpdateById(*dtos.EquipmentUpdate, *dtos.EquipmentFilter) (bool, error)
	UpdateByIds(*dtos.EquipmentUpdate, *dtos.EquipmentFilter) ([]string, error)
	UpdateByConditions(*dtos.EquipmentUpdate, *dtos.EquipmentFilter) ([]string, error)
	DeleteById(*dtos.EquipmentFilter) (bool, error)
	DeleteByIds(*dtos.EquipmentFilter) ([]string, error)
	DeleteByConditions(*dtos.EquipmentFilter) ([]string, error)
	FindById(*dtos.EquipmentFilter) (*dtos.EquipmentGet, error)
	FindByIds(*dtos.EquipmentFilter) ([]dtos.EquipmentGet, error)
	FindByConditions(*dtos.EquipmentFilter) ([]dtos.EquipmentGet, error)
}

type Equipment struct {
	repository EquipmentRepository
}

func NewEquipment(repository EquipmentRepository) Equipment {
	return Equipment{repository: repository}
}

func (service Equipment) Close() {
	/*service.repository.Close() // How to avoid including method Close in interface? */
}

func (service Equipment) CreateOne(equipmentCreate *dtos.EquipmentCreate) (string, error) {
	return service.repository.CreateOne(equipmentCreate)
}

func (service Equipment) CreateMany(equipmentCreate []dtos.EquipmentCreate) ([]string, error) {
	return service.repository.CreateMany(equipmentCreate)
}

func (service Equipment) UpdateById(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) (bool, error) {
	return service.repository.UpdateById(equipmentUpdate, equipmentFilter)
}

func (service Equipment) UpdateByIds(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	return service.repository.UpdateByIds(equipmentUpdate, equipmentFilter)
}

func (service Equipment) UpdateByConditions(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	return service.repository.UpdateByConditions(equipmentUpdate, equipmentFilter)
}

func (service Equipment) DeleteById(equipmentFilter *dtos.EquipmentFilter) (bool, error) {
	return service.repository.DeleteById(equipmentFilter)
}

func (service Equipment) DeleteByIds(equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	return service.repository.DeleteByIds(equipmentFilter)
}

func (service Equipment) DeleteByConditions(equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	return service.repository.DeleteByConditions(equipmentFilter)
}

func (service Equipment) FindById(equipmentFilter *dtos.EquipmentFilter) (*dtos.EquipmentGet, error) {
	return service.repository.FindById(equipmentFilter)
}

func (service Equipment) FindByIds(equipmentFilter *dtos.EquipmentFilter) ([]dtos.EquipmentGet, error) {
	return service.repository.FindByIds(equipmentFilter)
}

func (service Equipment) FindByConditions(equipmentFilter *dtos.EquipmentFilter) ([]dtos.EquipmentGet, error) {
	return service.repository.FindByConditions(equipmentFilter)
}
