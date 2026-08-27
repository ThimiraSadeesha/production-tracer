package schema

import (
	"time"

	"gorm.io/gorm"
)

type Machine struct {
	ID              int            `json:"id" gorm:"primaryKey;column:id"`
	MachineCode     string         `json:"machineCode" gorm:"type:varchar(100);column:machine_code"`
	MachineName     string         `json:"machineName" gorm:"type:varchar(255);column:machine_name"`
	Description     string         `json:"description" gorm:"type:varchar(255);column:description"`
	Status          string         `json:"status" gorm:"type:enum('Offline','Idle','Running','Paused','Stopped','Maintenance');not null;default:'Idle';column:status"`
	HourlyOutput    *float64       `json:"hourlyOutput" gorm:"column:hourly_output"`
	MachineTypeID   *int           `json:"machineTypeId" gorm:"column:machine_type_id"`
	MakeReadyTime   float64        `json:"makeReadyTime" gorm:"column:make_ready_time;default:0"`
	UnitOfMachineID *int           `json:"unitOfMachineId" gorm:"column:unit_of_machine_id"`
	CreatedBy       *string        `json:"createdBy" gorm:"type:varchar(255);column:created_by"`
	UpdatedBy       *string        `json:"updatedBy" gorm:"type:varchar(255);column:updated_by"`
	DeletedBy       *string        `json:"deletedBy" gorm:"type:varchar(255);column:deleted_by"`
	CreatedAt       time.Time      `json:"createdAt" gorm:"autoCreateTime;column:created_at"`
	UpdatedAt       time.Time      `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deletedAt" gorm:"index;column:deleted_at"`

	// Relationships
	MachineType  *MachineType  `json:"machineType,omitempty" gorm:"foreignKey:MachineTypeID"`
	Unit         *Unit         `json:"unit,omitempty" gorm:"foreignKey:UnitOfMachineID"`
	ShiftRosters []ShiftRoster `json:"shiftRosters,omitempty" gorm:"foreignKey:MachineID"`
}

func (Machine) TableName() string {
	return "tbl_machine"
}
