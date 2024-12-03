package repository

import (
	"fmt"
	"strings"
	"github.com/jmoiron/sqlx"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/dtos"
	"github.com/Melanjnk/equipment-monitor/internal/app/registry_service/model"
)

func parseRows(rows *sqlx.Rows, err error) ([]string, error) {
	if err == nil {
		var ids []string
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
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

func joinIds(ids []string) string {
	if i := len(ids); i == 0 {
		return ""
	} else {
		var builder strings.Builder
		for {
			i--
			builder.WriteByte(byte('\''))
			builder.WriteString(ids[i])
			builder.WriteByte(byte('\''))
			if i == 0 {
				return builder.String()
			}
			builder.WriteByte(byte(','))
		}
	}
}

func updateSQL(equipmentUpdate *dtos.EquipmentUpdate, base string, extra string) string {
	var set string
	if equipmentUpdate.Parameters == nil {
		if equipmentUpdate.Status == nil {
			return `` // Nothing to update
		}
		set = `status=:status`
	} else {
		if equipmentUpdate.Status == nil {
			set = `parameters=:parameters`
		} else {
			set = `status=:status,parameters=:parameters`
		}
	}
	return fmt.Sprintf(base, set, extra)
}

// processArgs checks whether `kind`/`no_kind` and `status`/`no_status` arguments are used,
// and serializes corresponding conditions to `conditions` slice.
func processArgs[A model.EquipmentKind|model.OperationalStatus](builder strings.Builder, conditions *[]string, argName string, arg, noArg []A) {
	var a []A
	if arg != nil {
		a = arg
		// Write `parameter=`
		builder.WriteString(argName)
		builder.WriteByte(byte('='))
	} else if noArg != nil {
		a = noArg
		// Write `parameter<>`
		builder.WriteString(argName)
		builder.WriteByte(byte('<'))
		builder.WriteByte(byte('>'))
	} else {
		return
	}
	switch i := len(a); i {
		case 0:
			return
		case 1:
			// Write `value`
			fmt.Fprintf(&builder, `%d`, a[0])
		default:
			// Write `ANY(ARRAY(value0, value1, ...))`
			fmt.Fprintf(&builder, ` ANY(ARRAY[%d`, a[0])
			for i > 1 {
				i--
				fmt.Fprintf(&builder, `,%d`, a[i])
			} 
			builder.WriteByte(byte(']'))
			builder.WriteByte(byte(')'))
	}
	*conditions = append(*conditions, builder.String())
	builder.Reset()
}

func getConditions(equipmentFilter *dtos.EquipmentFilter) string {
	conditions := make([]string, 0, 6) // Capacity == maximal number of possible simultaneous conditions
	var builder strings.Builder

	processArgs(builder, &conditions, "kind", equipmentFilter.Kinds, equipmentFilter.NoKinds)
	processArgs(builder, &conditions, "status", equipmentFilter.Statuses, equipmentFilter.NoStatuses)

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

	if i := len(conditions); i == 0 {
		return ""
	} else {
		// Write ` WHERE condition0` (` AND condition1` ...)
		builder.WriteString(" WHERE ")
		for {
			i--
			builder.WriteString(conditions[i])
			if i == 0 {
				break
			}
			builder.WriteString(" AND ")
		}
		return builder.String()
	}
}
