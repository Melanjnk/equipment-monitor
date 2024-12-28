package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/dtos"
)

type Equipment struct {
	db *sqlx.DB
	insertOne *sqlx.NamedStmt
	deleteOne *sqlx.NamedStmt
}

const (
	insertEquipment = `INSERT INTO equipment(kind, parameters) VALUES(:kind, :parameters) RETURNING id;`
	deleteOneEquipment = `DELETE FROM equipment WHERE id=:id RETURNING id;`
)

func NewEquipment(db *sqlx.DB) *Equipment {
	MustPrepareNamed := func(query string) *sqlx.NamedStmt {
		if statement, err := db.PrepareNamed(query); err != nil {
			panic(err)
		} else {
			return statement
		}
	}

	return &Equipment{
		db,
		MustPrepareNamed(insertEquipment),
		MustPrepareNamed(deleteOneEquipment),
	}
}

func (repository *Equipment) CreateOne(equipmentCreate *dtos.EquipmentCreate) (string, error) {
	var id string
	insertOne, err := repository.db.PrepareNamed(insertEquipment)
	if err != nil {
		return "", err
	}
	err = insertOne.Get(&id, equipmentCreate)
	return id, err
}

func (repository *Equipment) CreateMany(equipmentCreate []dtos.EquipmentCreate) ([]string, error) {
	return parseRows(repository.db.NamedQuery(insertEquipment, equipmentCreate))
}

func (repository *Equipment) UpdateById(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) (bool, error) {
	if query, args := updateSQL(equipmentUpdate, equipmentFilter); len(query) > 0 {
		updateOne, err := repository.db.PrepareNamed(query)
		if err == nil {
			var updatedId string
			err = updateOne.Get(&updatedId, args)
			if err == nil {
				return true, nil
			}
			if !errors.Is(err, sql.ErrNoRows) {
				return false, err
			}
		}
	}
	return false, nil
}

func (repository *Equipment) UpdateByIds(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	if query, args := updateSQL(equipmentUpdate, equipmentFilter); len(query) == 0 {
		return nil, nil
	} else {
		return parseRows(repository.db.NamedQuery(repository.db.Rebind(query), args))
	}
}

func (repository *Equipment) UpdateByConditions(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	if query, args := updateSQL(equipmentUpdate, equipmentFilter); len(query) == 0 {
		return nil, nil
	} else {
		preparedQuery, err := repository.db.PrepareNamed(repository.db.Rebind(query))
		if err == nil {
			var ids []string
			err = preparedQuery.Select(&ids, args)
			if err == nil {
				return ids, nil
			}
		}
		return nil, err
	}
}

func (repository *Equipment) DeleteById(equipmentFilter *dtos.EquipmentFilter) (bool, error) {
	var deletedId string
	err := repository.deleteOne.Get(&deletedId, map[string]any{`id`: equipmentFilter.Ids[0]})
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, err
}

func (repository *Equipment) DeleteByIds(equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	if query, args, err := sqlx.In(`DELETE FROM equipment WHERE id IN(?) RETURNING id;`, equipmentFilter.Ids); err != nil {
		return nil, err
	} else {
		return parseRows(repository.db.Queryx(repository.db.Rebind(query), args...))
	}
}

func (repository *Equipment) DeleteByConditions(equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	query, args := deleteSQL(equipmentFilter)
	preparedQuery, err := repository.db.PrepareNamed(repository.db.Rebind(query))
	if err == nil {
		var ids []string
		err = preparedQuery.Select(&ids, args)
		if err == nil {
			return ids, nil
		}
	}
	return nil, err
}

func (repository *Equipment) FindById(equipmentFilter *dtos.EquipmentFilter) (*dtos.EquipmentGet, error) {
	query, args := findSQL(equipmentFilter)
	preparedQuery, err := repository.db.PrepareNamed(query); 
	if err == nil {
		var equipment dtos.EquipmentGet
		err = preparedQuery.Get(&equipment, args)
		if err == nil {
			return &equipment, nil
		}
	}
	return nil, err
}

func (repository *Equipment) FindByIds(equipmentFilter *dtos.EquipmentFilter) ([]dtos.EquipmentGet, error) {
	query, args, err := sqlx.In(`SELECT id, kind, status, parameters, created_at, updated_at FROM equipment WHERE id IN(?)`, equipmentFilter.Ids)
	if err == nil {
		var equipments []dtos.EquipmentGet
		if err = repository.db.Select(&equipments, repository.db.Rebind(query), args...); err == nil {
			return equipments, nil
		}
	}
	return nil, err
}

func (repository *Equipment) FindByConditions(equipmentFilter *dtos.EquipmentFilter) ([]dtos.EquipmentGet, error) {
	query, args, err := sqlx.Named(findSQL(equipmentFilter))
	if err == nil {
		query, args, err = sqlx.In(query, args...)
		if err == nil {
			query = repository.db.Rebind(query)
			var preparedQuery *sqlx.Stmt
			preparedQuery, err = repository.db.Preparex(query)
			if err == nil {
				var equipments []dtos.EquipmentGet
				err = preparedQuery.Select(&equipments, args...)
				if err == nil {
					return equipments, nil
				}
			}
		}
	}
	return nil, err
}
