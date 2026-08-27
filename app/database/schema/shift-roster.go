package schema

import (
	"time"

	"gorm.io/gorm"
)

type ShiftRoster struct {
	ID         int            `json:"id" gorm:"primaryKey;column:id"`
	Date       time.Time      `json:"date" gorm:"column:date"`
	Status     string         `json:"status" gorm:"type:varchar(100);column:status"`
	OperatorID int            `json:"operatorId" gorm:"column:operator_id"`
	ShiftID    int            `json:"shiftId" gorm:"column:shift_id"`
	MachineID  int            `json:"machineId" gorm:"column:machine_id"`
	ProcessID  *int           `json:"processId,omitempty" gorm:"column:process_id"`
	CreatedBy  *string        `json:"createdBy" gorm:"type:varchar(255);column:created_by"`
	UpdatedBy  *string        `json:"updatedBy" gorm:"type:varchar(255);column:updated_by"`
	DeletedBy  *string        `json:"deletedBy" gorm:"type:varchar(255);column:deleted_by"`
	CreatedAt  time.Time      `json:"createdAt" gorm:"autoCreateTime;column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"index;column:deleted_at"`

	// Relationships
	Operator *Operator `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
	Shift    *Shift    `json:"shift,omitempty" gorm:"foreignKey:ShiftID"`
	Machine  *Machine  `json:"machine,omitempty" gorm:"foreignKey:MachineID"`
	Process  *Process  `json:"process,omitempty" gorm:"foreignKey:ProcessID"`
}

func (ShiftRoster) TableName() string {
	return "tbl_shift_roster"
}
