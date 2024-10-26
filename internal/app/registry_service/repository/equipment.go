package repository

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/model"
)

type Equipment struct {
	db *sqlx.DB
}

func NewEquipment(db *sqlx.DB) Equipment {
	return Equipment{db: db}
}

func (eqr *Equipment) List() ([]model.Equipment, error) {
	var eql []model.Equipment
	err := eqr.db.Select(&eql, `SELECT id, type, status, parameters, created_at, updated_at FROM equipment`)
	if err != nil {
		return nil, err
	}
	return eql, nil
}

func (eqr *Equipment) Create(et model.EquipmentType, ep model.Params) (uuid.UUID, error) {
	// Marshal the Parameters map to JSON
	jep, err := json.Marshal(ep)
	if err != nil {
		return uuid.UUID{}, err
	}

	var id uuid.UUID
	for {
		// Generate a UUID version 6 (using a library):
		id, err = uuid.NewV6()
		if err != nil {
			return uuid.UUID{}, err
		}
		// Verify that id does not exist in the database:
		eq, err := eqr.FindById(id)
		if err != nil {
			return uuid.UUID{}, err
		}
		if eq == nil {
			break
		}
	}
	// Insert into the database
	_, err = eqr.db.NamedExec(
		`INSERT INTO equipment (id, type, status, parameters) VALUES (:id, :type, :status, :parameters)`,
		map[string]interface{}{
			"id":         id,
			"type":       et,
			"status":     model.Operational,
			"parameters": string(jep),
		},
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (eqr *Equipment) Update(id uuid.UUID, es *model.OperationalStatus, ep model.Params) error {
	var set, p string
	if ep == nil {
		if es == nil {
			return nil // Nothing to update
		}
		set = "status=:status"
		p = ""
	} else {
		jep, err := json.Marshal(ep)
		if err != nil {
			return err
		}
		if es == nil {
			set = "parameters=:parameters, updated_at=:updated_at"
		} else {
			set = "status=:status, parameters=:parameters, updated_at=:updated_at"
		}
		p = string(jep)
	}
	_, err := eqr.db.NamedExec(
		fmt.Sprintf("UPDATE equipment SET %s WHERE id=:id", set),
		map[string]interface{}{
			"id":			id,
			"status": 		es,
			"parameters":	p,
			"updated_at":	time.Now(),
		},
	)
	return err
}

func (eqr *Equipment) FindById(id uuid.UUID) (*model.Equipment, error) {
	var eq model.Equipment
	err := eqr.db.Get(&eq, `SELECT id, type, status, parameters FROM equipment WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	return &eq, nil
}

func (eqr *Equipment) RemoveById(id uuid.UUID) (bool, error) {
	_, err := eqr.db.Exec(`DELETE FROM equipment WHERE id=$1`, id)
	return err == nil, err
}
