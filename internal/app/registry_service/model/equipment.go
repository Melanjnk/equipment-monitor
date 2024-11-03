package model

import (
	"strings"
	"time"
	"github.com/gofrs/uuid"
)

type EquipmentKind int16
const (
	CNCMachine EquipmentKind = iota
	ConveyorBelt
	DrillMachine
	RoboticArm
)

func (kind EquipmentKind) IsValid() bool {
	return kind >= CNCMachine && kind <= RoboticArm
}

func (kind EquipmentKind) String() string {
	switch kind {
		case CNCMachine:	return "CNCMachine"
		case ConveyorBelt:	return "ConveyorBelt"
		case DrillMachine:	return "DrillMachine"
		case RoboticArm:	return "RoboticArm"
	}
	return "";
} 


func ParseEquipmentKind(str string) *EquipmentKind {
	var kind EquipmentKind
	switch strings.ToLower(str) {
		case "cncmachine":		kind = CNCMachine
		case "conveyorbelt":	kind = ConveyorBelt
		case "drillmachine":	kind = DrillMachine
		case "roboticarm":		kind = RoboticArm
		default:				return nil
	}
	return &kind
}


type OperationalStatus int8
const (
	Operational OperationalStatus = iota
	UnderMaintenance
	Decommissioned
)

func (status OperationalStatus) IsValid() bool {
	return status >= Operational && status <= Decommissioned
}

func (status OperationalStatus) String() string {
	switch status {
		case Operational:		return "Operational"
		case UnderMaintenance:	return "UnderMaintenance"
		case Decommissioned:	return "Decommissioned"
	}
	return ""
}

func ParseOperationalStatus(str string) *OperationalStatus {
	var status OperationalStatus
	switch strings.ToLower(str) {
		case "operational":			status = Operational
		case "undermaintenance":	status = UnderMaintenance
		case "decommissioned":		status = Decommissioned
		default:					return nil
	}
	return &status
}


type Equipment struct {
	Id			uuid.UUID			`db:"id,pk"								json:"equipment_id"`
	Kind		EquipmentKind		`db:"kind,not null,type:smallserial"	json:"kind"`
	Status		OperationalStatus	`db:"status,not null,type:smallserial"	json:"status"`
	Parameters	[]byte				`db:"parameters,not null,type:jsonb"	json:"parameters"`
	CreatedAt	time.Time			`db:"created_at,not null"				json:"created_at"`
	UpdatedAt	time.Time			`db:"updated_at,not null"				json:"updated_at"`
}

func NewEquipment(kind EquipmentKind, parameters []byte) Equipment {
	return Equipment {
		Kind:		kind,
		Status:		Operational,	// Default status on creation
		Parameters:	parameters,
	}
}
