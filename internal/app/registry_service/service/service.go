package service

import (
	"github.com/gofrs/uuid"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/model"
)

type EquipmentRepo interface {
	List() ([]model.Equipment, error)
	Create(et model.EquipmentType, ep model.Params) (uuid.UUID, error)
	Update(id uuid.UUID, status *model.OperationalStatus, parameters model.Params) error
	FindById(Id uuid.UUID) (*model.Equipment, error)
	RemoveById(Id uuid.UUID) (bool, error)
}

type Equipment struct {
	repo EquipmentRepo
}

func NewEquipment(repo EquipmentRepo) Equipment {
	return Equipment{repo: repo}
}

func (eqs *Equipment) List() ([]model.Equipment, error) {
	return eqs.repo.List()
}

func (eqs *Equipment) Create(eqt model.EquipmentType, eqp model.Params) (uuid.UUID, error) {
	return eqs.repo.Create(eqt, eqp)
}

func (eqs *Equipment) Update(id uuid.UUID, eqos *model.OperationalStatus, eqp model.Params) error {
	return eqs.repo.Update(id, eqos, eqp)
}

func (eqs *Equipment) Get(eqId string) (*model.Equipment, error) {
	id, err := uuid.FromString(eqId)
	if err != nil {
		return nil, err
	}

	return eqs.repo.FindById(id)
}

func (eqs *Equipment) Delete(eqId string) (bool, error) {
	id, err := uuid.FromString(eqId)
	if err != nil {
		return false, err
	}
	return eqs.repo.RemoveById(id)
}
