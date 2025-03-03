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
	insertEquipment = `INSERT INTO equipment(kind,parameters) VALUES(:kind,:parameters) RETURNING id;`
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

func (repository *Equipment) Close() {
	repository.insertOne.Close()
	repository.deleteOne.Close()
}

func (repository *Equipment) CreateOne(equipmentCreate *dtos.EquipmentCreate) (string, error) {
	var id string
	err := repository.insertOne.Get(&id, equipmentCreate)
	return id, err
}

func (repository *Equipment) CreateMany(equipmentCreate []dtos.EquipmentCreate) ([]string, error) {
	return parseRows(repository.db.NamedQuery(insertEquipment, equipmentCreate))
}

func (repository *Equipment) UpdateById(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) (bool, error) {
	query, args, err := updateSQL(equipmentUpdate, equipmentFilter)
	if err == nil {
		if len(query) > 0 {
			var updateOne *sqlx.NamedStmt
			if updateOne, err = repository.db.PrepareNamed(query); err == nil {
				defer updateOne.Close()

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
	}
	return false, err
}

func (repository *Equipment) UpdateByIds(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	if query, args, err := updateSQL(equipmentUpdate, equipmentFilter); err != nil {
		return nil, err
	} else if len(query) == 0 {
		return nil, nil
	} else {
		preparedQuery, err := repository.db.Preparex(repository.db.Rebind(query))
		if err != nil {
			return nil, err
		}
		var ids []string
		if argSlice, ok := args.([]any); ok {
			err = preparedQuery.Select(&ids, argSlice...)
		} else {
			err = preparedQuery.Select(&ids, args)
		}
		return ids, err
	}
}

func (repository *Equipment) UpdateByConditions(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) ([]string, error) {
	query, args, err := updateSQL(equipmentUpdate, equipmentFilter)
	if err == nil {
		var updateByConditions *sqlx.NamedStmt
		if updateByConditions, err = repository.db.PrepareNamed(repository.db.Rebind(query)); err == nil {
			defer updateByConditions.Close()

			var ids []string
			if err = updateByConditions.Select(&ids, args); err == nil {
				if ids != nil {
					return ids, nil
				}
				return make([]string, 0, 0), nil
			}
		}
	}
	return nil, err
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
	query, args, err := deleteSQL(equipmentFilter)
	if err == nil {
		var deleteByConditions *sqlx.NamedStmt
		deleteByConditions, err = repository.db.PrepareNamed(repository.db.Rebind(query))
		if err == nil {
			defer deleteByConditions.Close()

			var ids []string
			err = deleteByConditions.Select(&ids, args)
			if err == nil {
				if ids != nil {
					return ids, nil
				}
				return make([]string, 0, 0), nil
			}
		}
	}
	return nil, err
}

func (repository *Equipment) FindById(equipmentFilter *dtos.EquipmentFilter) (*dtos.EquipmentGet, error) {
	query, args, err := findSQL(equipmentFilter)
	if err == nil {
		var findById *sqlx.NamedStmt
		findById, err = repository.db.PrepareNamed(query); 
		if err == nil {
			defer findById.Close()

			var equipment dtos.EquipmentGet
			err = findById.Get(&equipment, args)
			if err == nil {
				return &equipment, nil
			}
		}
	}
	return nil, err
}

func (repository *Equipment) FindByIds(equipmentFilter *dtos.EquipmentFilter) ([]dtos.EquipmentGet, error) {
	query, args, err := sqlx.In(`SELECT id,kind,status,parameters,created_at,updated_at FROM equipment WHERE id IN(?)`, equipmentFilter.Ids)
	if err == nil {
		var equipments []dtos.EquipmentGet
		if err = repository.db.Select(&equipments, repository.db.Rebind(query), args...); err == nil {
			return equipments, nil
		}
	}
	return nil, err
}

func (repository *Equipment) FindByConditions(equipmentFilter *dtos.EquipmentFilter) ([]dtos.EquipmentGet, error) {
	query, args, err := findSQL(equipmentFilter)
	if err == nil {
		if query, args, err = sqlx.Named(query, args); err == nil {
			var findByConditions *sqlx.Stmt
			if findByConditions, err = repository.db.Preparex(repository.db.Rebind(query)); err == nil {
				defer findByConditions.Close()

				var equipments []dtos.EquipmentGet
				if argSlice, ok := args.([]any); ok {
					err = findByConditions.Select(&equipments, argSlice...)
				} else {
					err = findByConditions.Select(&equipments, args)
				}
				if err == nil {
					if equipments != nil {
						return equipments, nil
					}
					// Nothing is found; avoid returning nil
					return make([]dtos.EquipmentGet, 0, 0), nil
				}
			}
		}
	}
	return nil, err
}
