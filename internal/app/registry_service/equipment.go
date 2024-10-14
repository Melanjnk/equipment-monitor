package registry_service

import (
	"github.com/gofrs/uuid"
	"github.com/google/uuid"
)

type EquipmentType int16
const(
	CNCMachine = iota,
	ConveyorBelt,
	DrillMachine,
	RoboticArm,
)

type OperationalStatus int8
const (
	Operational OperationalStatus = iota,
	UnderMaintenance,
	Decommissioned,
)

type Equipment struct {
	Id		uuid.UUID		`pg:"id,pk"		json:"equipment_id"`
	Type		int16			`pg:"type"		json:"type"`		// pg: smallint
	Status		OperationalStatus	`pg:"status"		json:"status"`
	Parameters	map[string]interface{}	`pg:"parameters"	json:"parameters"`	// pg: jsonb
}
