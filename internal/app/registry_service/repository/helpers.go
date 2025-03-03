package repository

import (
	"fmt"
	"strings"
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/dtos"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/model"
)

type arguments = map[string]any

func parseRows(rows *sqlx.Rows, err error) ([]string, error) {
	if err == nil {
		var ids []string
		for rows.Next() {
			var id string
			err = rows.Scan(&id)
			if err != nil {
				rows.Close()
				goto ERROR
			}
			ids = append(ids, id)
		}
		return ids, nil
	}
ERROR:
	return nil, err
}

// processParam checks whether `id`/`no_id`, `kind`/`no_kind` and `status`/`no_status` arguments are used,
// and serializes corresponding conditions to `conditions` slice.
func processParam[P model.EquipmentKind|model.OperationalStatus|string](conditions *[]string, args arguments, paramName string, param, noParam []P) bool {
	var p []P
	var b bool
	var builder strings.Builder
	if param != nil {
		p = param
		b = true
	} else if noParam != nil {
		p = noParam
		b = false
	} else {
		return false
	}
	builder.WriteString(paramName)
	switch i := len(p); i {
		case 0:
			return false
		case 1:
			if b {
				builder.WriteByte('=')
			} else {
				builder.WriteByte('<')
				builder.WriteByte('>')
			}
			builder.WriteByte(':')
			builder.WriteString(paramName)
			args[paramName] = p[0]
			b = false
		default:
			if !b {
				builder.WriteString(` NOT`)
			}
			builder.WriteString(` IN(:`)
			builder.WriteString(paramName)
			builder.WriteByte(')')
			args[paramName] = p
			b = true
	}
	*conditions = append(*conditions, builder.String())
	return b
}

// Return true if IN operator is used at least once, false otherwise.
func getConditions(builder *strings.Builder, equipmentFilter *dtos.EquipmentFilter, args arguments) bool {
	conditions := make([]string, 0, 7) // Capacity == maximal number of possible simultaneous conditions

	result := processParam(&conditions, args, `id`, equipmentFilter.Ids, equipmentFilter.NoIds)
	result = processParam(&conditions, args, `kind`, equipmentFilter.Kinds, equipmentFilter.NoKinds) || result
	result = processParam(&conditions, args, `status`, equipmentFilter.Statuses, equipmentFilter.NoStatuses) || result

	processTimeParam := func(param *time.Time, paramName string, condition string) {
		if param != nil {
			conditions = append(conditions, condition + paramName)
			args[paramName] = param
		}
	}

	processTimeParam(equipmentFilter.CreatedSince, `created_since`, `created_at>=:`)
	processTimeParam(equipmentFilter.CreatedUntil, `created_until`, `created_at<=:`)
	processTimeParam(equipmentFilter.UpdatedSince, `updated_since`, `updated_at>=:`)
	processTimeParam(equipmentFilter.UpdatedUntil, `updated_until`, `updated_at<=:`)

	if i := len(conditions); i > 0 {
		// Write ` WHERE condition0` (` AND condition1` ...)
		builder.WriteString(` WHERE `)
		for {
			i--
			builder.WriteString(conditions[i])
			if i == 0 {
				break
			}
			builder.WriteString(` AND `)
		}
	}
	return result
}

func updateSQL(equipmentUpdate *dtos.EquipmentUpdate, equipmentFilter *dtos.EquipmentFilter) (string, any, error) {
	const update = `UPDATE equipment SET `
	var args arguments
	var builder strings.Builder
	if equipmentUpdate.Parameters == nil {
		if equipmentUpdate.Status == nil {
			return ``, nil, nil // Nothing to update
		}
		builder.WriteString(update)
		builder.WriteString(`status=:new_status`)
		args = make(arguments)
		args[`new_status`] = equipmentUpdate.Status
	} else {
		args = make(arguments)
		args[`new_parameters`] = equipmentUpdate.Parameters
		builder.WriteString(update)
		if equipmentUpdate.Status == nil {
			builder.WriteString(`parameters=:new_parameters`)
		} else {
			builder.WriteString(`status=:new_status,parameters=:new_parameters`)
			args[`new_status`] = equipmentUpdate.Status
		}
	}
	hasIn := getConditions(&builder, equipmentFilter, args)
	builder.WriteString(` RETURNING id;`)
	return checkIn(hasIn, &builder, args)
}

func deleteSQL(equipmentFilter *dtos.EquipmentFilter) (string, any, error) {
	args := make(arguments)
	var builder strings.Builder
	builder.WriteString(`DELETE FROM equipment`)
	b := getConditions(&builder, equipmentFilter, args)
	builder.WriteString(` RETURNING id;`)
	return checkIn(b, &builder, args)
}

func findSQL(equipmentFilter *dtos.EquipmentFilter) (string, any, error) {
	args := make(arguments)
	var builder strings.Builder
	builder.WriteString(`SELECT id,kind,status,parameters,created_at,updated_at FROM equipment`)
	hasIn := getConditions(&builder, equipmentFilter, args)
	if l := len(equipmentFilter.Sort); l > 0 {
		builder.WriteString(` ORDER BY `)
		for i := 0; ; {
			builder.WriteString(equipmentFilter.Sort[i])
			if equipmentFilter.SortMask & (1 << i) != 0 {
				builder.WriteString(` DESC`)
			}
			i++
			if i == l {
				break
			}
			builder.WriteByte(',')
		}
	}
	if equipmentFilter.Limit != nil {
		fmt.Fprintf(&builder, ` LIMIT %d`, *equipmentFilter.Limit)
	}
	if equipmentFilter.Offset != nil {
		fmt.Fprintf(&builder, ` OFFSET %d`, *equipmentFilter.Offset)
	}
	builder.WriteByte(';')
	return checkIn(hasIn, &builder, args)
}

func checkIn(hasIn bool, builder *strings.Builder, args arguments) (string, any, error) {
	if hasIn {
		query, args, err := sqlx.Named(builder.String(), args)
		if err != nil {
			return ``, nil, err
		}
		return sqlx.In(query, args...)
	}
	return builder.String(), args, nil
}
