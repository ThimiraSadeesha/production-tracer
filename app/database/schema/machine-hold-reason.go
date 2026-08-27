package schema

import (
	"time"
)

type MachineHoldReason struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	MachineID    int       `json:"machine_id"`
	HoldReasonID int       `json:"hold_reason_id"`
	CreatedAt    time.Time `json:"created_at"`

	// Relations
	Machine    Machine    `json:"machine" gorm:"foreignKey:MachineID"`
	HoldReason HoldReason `json:"hold_reason" gorm:"foreignKey:HoldReasonID"`
}

func (MachineHoldReason) TableName() string {
	return "tbl_machine_hold_reason"
}
