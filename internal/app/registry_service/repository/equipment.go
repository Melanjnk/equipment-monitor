package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/dtos"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/model"
)

type Equipment struct {
	db *sqlx.DB
}

func NewEquipment(db *sqlx.DB) Equipment {
	return Equipment{db: db}
}

func (repository *Equipment) List() ([]*dtos.EquipmentGet, error) {
	var equipmentModels []model.Equipment
	err := repository.db.Select(&equipmentModels, `SELECT id, kind, status, parameters, created_at, updated_at FROM equipment`)
	if err != nil {
		return nil, err
	}
	equipmentGets := make([]*dtos.EquipmentGet, 0, len(equipmentModels))
	for _, equipmentModel := range equipmentModels {
		equipmentGets = append(equipmentGets, dtos.EquipmentGetFromModel(equipmentModel))
	}
	return equipmentGets, nil
}

func (repository *Equipment) Create(equipmentCreate *dtos.EquipmentCreate) (uuid.UUID, error) {
	for jsonifiedParameters, _ := json.Marshal(equipmentCreate.Parameters); ; {
		// Generate a UUID version 6 (using a library):
		id, err := uuid.NewV6()
		if err == nil {
			_, err = repository.db.NamedExec(
				`INSERT INTO equipment (id, kind, status, parameters) VALUES (:id, :kind, :status, :parameters)`,
				map[string]interface{}{
					"id":         id,
					"kind":       equipmentCreate.Kind,
					"status":     model.Operational,
					"parameters": jsonifiedParameters,
				},
			)
			if err == nil {
				return id, nil
			}
			if false {// TODO: Check id collision error
				continue
			}
		}
		return uuid.UUID{}, err
	}
}

func (repository *Equipment) Update(equipmentUpdate *dtos.EquipmentUpdate) (bool, error) {
	var set string
	var jsonifiedParameters []byte
	if equipmentUpdate.Parameters == nil {
		if equipmentUpdate.Status == nil {
			return false, nil // Nothing to update
		}
		set = "status=:status"
	} else {
		if equipmentUpdate.Status == nil {
			set = "parameters=:parameters, updated_at=:updated_at"
		} else {
			set = "status=:status, parameters=:parameters, updated_at=:updated_at"
		}
		jsonifiedParameters, _ = json.Marshal(*equipmentUpdate.Parameters)
	}
	return checkAffect(repository.db.NamedExec(
		fmt.Sprintf("UPDATE equipment SET %s WHERE id=:id", set),
		map[string]interface{}{
			"id":			equipmentUpdate.Id,
			"status": 		equipmentUpdate.Status,
			"parameters":	jsonifiedParameters,
			"updated_at":	time.Now(),
		},
	))
}

func (repository *Equipment) FindById(id uuid.UUID) (*dtos.EquipmentGet, error) {
	var equipmentModel model.Equipment
	err := repository.db.Get(&equipmentModel, `SELECT id, kind, status, parameters, created_at, updated_at FROM equipment WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	var parameters map[string]interface{}
	if err := json.Unmarshal(equipmentModel.Parameters, &parameters); err != nil {
		return nil, err
	}
	return dtos.EquipmentGetFromModel(equipmentModel), nil
}

func (repository *Equipment) RemoveById(id uuid.UUID) (bool, error) {
	return checkAffect(repository.db.Exec(`DELETE FROM equipment WHERE id=$1`, id))
}

func checkAffect(result sql.Result, err error) (bool, error) {
	if err == nil {
		var count int64
		if count, err = result.RowsAffected(); err == nil {
			return count > 0, nil
		}
	}
	return false, err
}
