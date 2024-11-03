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

func (repo *Equipment) List() ([]*dtos.EquipmentGet, error) {
	var eqs []model.Equipment
	err := repo.db.Select(&eqs, `SELECT id, kind, status, parameters, created_at, updated_at FROM equipment`)
	if err != nil {
		return nil, err
	}
	count := len(eqs)
	eqgs := make([]*dtos.EquipmentGet, count, count)
	for i, eq := range eqs {
		eqgs[i] = dtos.EquipmentGetFromModel(eq)
	}
	return eqgs, nil
}

func (repo *Equipment) Create(eqc *dtos.EquipmentCreate) (uuid.UUID, error) {
	for jParameters, _ := json.Marshal(eqc.Parameters); ; {
		// Generate a UUID version 6 (using a library):
		id, err := uuid.NewV6()
		if err == nil {
			_, err = repo.db.NamedExec(
				`INSERT INTO equipment (id, kind, status, parameters) VALUES (:id, :kind, :status, :parameters)`,
				map[string]interface{}{
					"id":         id,
					"kind":       eqc.Kind,
					"status":     model.Operational,
					"parameters": jParameters,
				},
			)
			if err == nil {
				return id, nil
			}
			if false {// TODO: Check id duplicate error
				continue
			}
		}
		return uuid.UUID{}, err
	}
}

func (repo *Equipment) Update(equ *dtos.EquipmentUpdate) (bool, error) {
	var set string
	var parameters []byte
	if equ.Parameters == nil {
		if equ.Status == nil {
			return false, nil // Nothing to update
		}
		set = "status=:status"
	} else {
		if equ.Status == nil {
			set = "parameters=:parameters, updated_at=:updated_at"
		} else {
			set = "status=:status, parameters=:parameters, updated_at=:updated_at"
		}
		parameters, _ = json.Marshal(*equ.Parameters)
	}
	return checkAffect(repo.db.NamedExec(
		fmt.Sprintf("UPDATE equipment SET %s WHERE id=:id", set),
		map[string]interface{}{
			"id":			equ.Id,
			"status": 		equ.Status,
			"parameters":	parameters,
			"updated_at":	time.Now(),
		},
	))
}

func (repo *Equipment) FindById(id uuid.UUID) (*dtos.EquipmentGet, error) {
	var eq model.Equipment
	err := repo.db.Get(&eq, `SELECT id, kind, status, parameters, created_at, updated_at FROM equipment WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	var parameters map[string]interface{}
	if err := json.Unmarshal(eq.Parameters, &parameters); err != nil {
		return nil, err
	}
	return dtos.EquipmentGetFromModel(eq), nil
}

func (repo *Equipment) RemoveById(id uuid.UUID) (bool, error) {
	return checkAffect(repo.db.Exec(`DELETE FROM equipment WHERE id=$1`, id))
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
