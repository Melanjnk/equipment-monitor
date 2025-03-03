package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
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


type EquipmentParameters map[string]any

func (parameters *EquipmentParameters) Scan(value any) error {
	switch v := value.(type) {
		case []byte:
			return json.Unmarshal(v, parameters)
		case string:
			return json.Unmarshal([]byte(v), parameters)
	}
	return fmt.Errorf("Failed to unmarshal JSONB value: %v", value)
}

func (parameters EquipmentParameters) Value() (driver.Value, error) {
    bytes, err := json.Marshal(parameters)
    return string(bytes), err
}


type Equipment struct {
	Id			string				`json:"equipment_id" db:"id,pk"`
	Kind		EquipmentKind		`json:"kind" db:"kind,not null"`
	Status		OperationalStatus	`json:"status" db:"status,not null"`
	Parameters	EquipmentParameters	`json:"parameters" db:"parameters,not null,type:jsonb"`
	CreatedAt	time.Time			`json:"created_at" db:"created_at,not null"`
	UpdatedAt	time.Time			`json:"updated_at" db:"updated_at,not null"`
}

func NewEquipment(kind EquipmentKind, parameters EquipmentParameters) Equipment {
	return Equipment {
		Kind:		kind,
		Status:		Operational,	// Default status on creation
		Parameters:	parameters,
	}
}
