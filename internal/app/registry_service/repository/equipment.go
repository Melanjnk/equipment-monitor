package repository

import (
	"encoding/json"
	"fmt"
	"strings"
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

func (repository *Equipment) List(equipmentFilter *dtos.EquipmentFilter) ([]*dtos.EquipmentGet, error) {
	var equipmentModels []model.Equipment
	query := `SELECT id, kind, status, parameters, created_at, updated_at FROM equipment`
	conditions := make([]string, 0, 6)

	if equipmentFilter.Kinds != nil {
		switch len(equipmentFilter.Kinds) {
			case 0:
				break
			case 1:
				conditions = append(conditions, "kind=" + string(figure(equipmentFilter.Kinds[0])))
			default:
				conditions = append(conditions, fmt.Sprintf("kind IN (%s)", joinIntegralArray(equipmentFilter.Kinds)))
		}
	} else if equipmentFilter.NoKinds != nil {
		switch len(equipmentFilter.NoKinds) {
			case 0:
				break
			case 1:
				conditions = append(conditions, "kind<>" + string(figure(equipmentFilter.NoKinds[0])))
			default:
				conditions = append(conditions, fmt.Sprintf("kind NOT IN (%s)", joinIntegralArray(equipmentFilter.NoKinds)))
		}
	}

	if equipmentFilter.Statuses != nil {
		switch len(equipmentFilter.Statuses) {
			case 0:
				break
			case 1:
				conditions = append(conditions, "status=" + string(figure(equipmentFilter.Statuses[0])))
			default:
				conditions = append(conditions, fmt.Sprintf("status IN (%s)", joinIntegralArray(equipmentFilter.Statuses)))
		}
	} else if equipmentFilter.NoStatuses != nil {
		switch len(equipmentFilter.NoStatuses) {
			case 0:
				break
			case 1:
				conditions = append(conditions, "status<>" + string(figure(equipmentFilter.NoStatuses[0])))
			default:
				conditions = append(conditions, fmt.Sprintf("status NOT IN (%s)", joinIntegralArray(equipmentFilter.NoStatuses)))
		}
	}

	if equipmentFilter.CreatedSince != nil {
		conditions = append(conditions, "created_at>=:created_since")
	}
	if equipmentFilter.CreatedUntil != nil {
		conditions = append(conditions, "created_at<=:created_until")
	}
	if equipmentFilter.UpdatedSince != nil {
		conditions = append(conditions, "updated_at>=:updated_since")
	}
	if equipmentFilter.UpdatedUntil != nil {
		conditions = append(conditions, "updated_at<=:updated_until")
	}
	var err error
	if len(conditions) == 0 {
		err = repository.db.Select(&equipmentModels, query)
	} else {
		preparedQuery, _ := repository.db.PrepareNamed(query + " WHERE " + strings.Join(conditions, " AND "))
		err = preparedQuery.Select(
			&equipmentModels,
			map[string]interface{}{
				"created_since":	equipmentFilter.CreatedSince,
				"created_until":	equipmentFilter.CreatedUntil,
				"updated_since":	equipmentFilter.UpdatedSince,
				"updated_until":	equipmentFilter.UpdatedUntil,
			},
		)
	}
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
			set = "parameters=:parameters"
		} else {
			set = "status=:status, parameters=:parameters"
		}
		jsonifiedParameters, _ = json.Marshal(*equipmentUpdate.Parameters)
	}
	return checkAffect(repository.db.NamedExec(
		fmt.Sprintf("UPDATE equipment SET %s, updated_at=:updated_at WHERE id=:id", set),
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
