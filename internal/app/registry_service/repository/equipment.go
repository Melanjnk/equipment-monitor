package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/dtos"
)

type Equipment struct {
	db *sqlx.DB
	createOne *sqlx.NamedStmt
}

const (
	insertEquipment = `INSERT INTO equipment (kind, parameters) VALUES (:kind, :parameters) RETURNING id`
)

func NewEquipment(db *sqlx.DB) *Equipment {
	createOne, _ := db.PrepareNamed(insertEquipment)
	return &Equipment{
		db: db,
		createOne: createOne,
	}
}

func (repository *Equipment) CreateOne(equipmentCreate *dtos.EquipmentCreate) (string, error) {
	var id string
	err := repository.createOne.Get(&id, equipmentCreate)
	return id, err
}

func (repository *Equipment) CreateMany(equipmentCreate []dtos.EquipmentCreate) ([]string, error) {
	return parseRows(repository.db.NamedQuery(insertEquipment, equipmentCreate))
}

func (repository *Equipment) UpdateById(equipmentUpdate *dtos.EquipmentUpdate, id string) error {
	if sql := updateSQL(equipmentUpdate, `UPDATE equipment SET %s WHERE id=%s`, id); len(sql) == 0 {
		return nil
	} else {
		_, err := repository.db.NamedExec(sql, equipmentUpdate)
		return err
	}
}

func (repository *Equipment) UpdateByIds(equipmentUpdate *dtos.EquipmentUpdate, ids []string) ([]string, error) {
	if sql := updateSQL(equipmentUpdate, `UPDATE equipment SET %s WHERE id IN(%s) RETURNING id`, joinIds(ids)); len(sql) == 0 {
		return nil, nil
	} else {
		return parseRows(repository.db.NamedQuery(sql, equipmentUpdate))
	}
}

func (repository *Equipment) UpdateByConditions(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	if sql := updateSQL(equipmentUpdate, `UPDATE equipment SET %s %s RETURNING id`, getConditions(equipmentFilter)); len(sql) == 0 {
		return nil, nil
	} else {
		preparedQuery, err := repository.db.PrepareNamed(sql)
		if err == nil {
			var ids []string
			err = preparedQuery.Select(
				&ids,
				equipmentUpdate,
				// TODO: unite `equipmentUpdate` fields and commented map below:
				/*map[string]interface{}{
					"created_since":	equipmentFilter.CreatedSince,
					"created_until":	equipmentFilter.CreatedUntil,
					"updated_since":	equipmentFilter.UpdatedSince,
					"updated_until":	equipmentFilter.UpdatedUntil,
				},*/
			)
			if err == nil {
				return ids, nil
			}
		}
		return nil, err
	}
}

func (repository *Equipment) DeleteById(id string) error {
	var deletedId string
	err := repository.db.Get(&deletedId, fmt.Sprintf(`DELETE FROM equipment WHERE id='%s' RETURNING id`, id))
	if err == nil {
		if id != deletedId {
			err = fmt.Errorf("Unable to find equipment #%s for deleting", id)
		}
	}
	return err
}

func (repository *Equipment) DeleteByIds(ids []string) ([]string, error) {
	return parseRows(repository.db.Queryx(
		fmt.Sprintf(`DELETE FROM equipment WHERE id IN(%s) RETURNING id`, joinIds(ids)),
	))
}

func (repository *Equipment) DeleteByConditions(equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	preparedQuery, err := repository.db.PrepareNamed(fmt.Sprintf(
		`DELETE FROM equipment%s RETURNING id`,
		getConditions(equipmentFilter),
	))
	if err == nil {
		var ids []string
		err = preparedQuery.Select(
			&ids,
			map[string]interface{}{
				"created_since":	equipmentFilter.CreatedSince,
				"created_until":	equipmentFilter.CreatedUntil,
				"updated_since":	equipmentFilter.UpdatedSince,
				"updated_until":	equipmentFilter.UpdatedUntil,
			},
		)
		if err == nil {
			return ids, nil
		}
	}
	return nil, err
}

func (repository *Equipment) FindById(id string) (*dtos.EquipmentGet, error) {
	var equipment dtos.EquipmentGet
	if err := repository.db.Get(
		&equipment,
		`SELECT id, kind, status, parameters, created_at, updated_at FROM equipment WHERE id=$1`,
		id,
	); err != nil {
		return nil, err
	}
	return &equipment, nil
}

func (repository *Equipment) FindByIds(ids []string) ([]dtos.EquipmentGet, error) {
	var equipments []dtos.EquipmentGet
	if err := repository.db.Select(
		&equipments,
		fmt.Sprintf(`SELECT id, kind, status, parameters, created_at, updated_at FROM equipment WHERE id IN(%s)`, joinIds(ids)),
	); err != nil {
		return nil, err
	}
	return equipments, nil
}

func (repository *Equipment) FindByConditions(equipmentFilter *dtos.EquipmentFilter) ([]dtos.EquipmentGet, error) {
	preparedQuery, err := repository.db.PrepareNamed(fmt.Sprintf(
		`SELECT id, kind, status, parameters, created_at, updated_at FROM equipment%s`,
		getConditions(equipmentFilter),
	))
	if err == nil {
		equipments := make([]dtos.EquipmentGet, 0, 0) // Dummy value to prevent returning nil for empty result
		err = preparedQuery.Select(
			&equipments,
			map[string]interface{}{
				"created_since":	equipmentFilter.CreatedSince,
				"created_until":	equipmentFilter.CreatedUntil,
				"updated_since":	equipmentFilter.UpdatedSince,
				"updated_until":	equipmentFilter.UpdatedUntil,
			},
		)
		if err == nil {
			return equipments, nil
		}
	}
	return nil, err
}
