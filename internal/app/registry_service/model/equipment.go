package model

import (
	"strings"
	"time"
	"github.com/gofrs/uuid"
)

type EquipmentType int16
const (
	CNCMachine EquipmentType = iota
	ConveyorBelt
	DrillMachine
	RoboticArm
)

func (et EquipmentType) String() string {
	switch et {
		case CNCMachine:	return "CNCMachine"
		case ConveyorBelt:	return "ConveyorBelt"
		case DrillMachine:	return "DrillMachine"
		case RoboticArm:	return "RoboticArm"
	}
	return "";
} 

func ParseEquipmentType(s string) *EquipmentType {
	var result EquipmentType
	switch strings.ToLower(s) {
		case "cncmachine":		result = CNCMachine
		case "conveyorbelt":	result = ConveyorBelt
		case "drillmachine":	result = DrillMachine
		case "roboticarm":		result = RoboticArm
		default:				return nil
	}
	return &result
}

type OperationalStatus int8
const (
	Operational OperationalStatus = iota
	UnderMaintenance
	Decommissioned
)

func (os OperationalStatus) String() string {
	switch os {
		case Operational:		return "Operational"
		case UnderMaintenance:	return "UnderMaintenance"
		case Decommissioned:	return "Decommissioned"
	}
	return ""
}

func ParseOperationalStatus(s string) *OperationalStatus {
	var result OperationalStatus
	switch strings.ToLower(s) {
		case "operational":			result = Operational
		case "undermaintenance":	result = UnderMaintenance
		case "decommissioned":		result = Decommissioned
		default:					return nil
	}
	return &result
}

type Params map[string]interface{}

type Equipment struct {
	Id			uuid.UUID			`db:"id,pk"								json:"equipment_id"`
	Type		EquipmentType		`db:"type,not null"						json:"type"`
	Status		OperationalStatus	`db:"status,not null"					json:"status"`
	Parameters	Params				`db:"parameters,not null,type:jsonb"	json:"parameters"`
	CreatedAt	time.Time			`db:"created_at,not null"				json:"created_at"`
	UpdatedAt	time.Time			`db:"updated_at,not null"				json:"updated_at"`
}

func NewEquipment(t EquipmentType, p Params) Equipment {
	return Equipment {
		Type:		t,
		Status:		Operational,	// Default status on creation
		Parameters:	p,
	}
}
