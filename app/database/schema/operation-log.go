package schema

import (
	"time"
)

type OperationLog struct {
	ID             int       `json:"id" gorm:"primaryKey;column:id"`
	PauseDateTime  time.Time `json:"pauseDateTime" gorm:"column:pause_date_time"`
	ResumeDateTime time.Time `json:"resumeDateTime" gorm:"column:resume_date_time"`
	JobOperationID int       `json:"jobOperationId" gorm:"column:job_operation_id"`
	MachineID      int       `json:"machineId" gorm:"column:machine_id"`
	PausedRosterID *int      `json:"pausedRosterId" gorm:"column:paused_roster_id"`
	ResumeRosterID *int      `json:"resumeRosterId" gorm:"column:resume_roster_id"`
	ReasonID       *int      `json:"reasonId" gorm:"column:reason"`
	Remark         *string   `json:"remark" gorm:"type:varchar(255);column:remark"`

	MachineHoldReason *MachineHoldReason          `json:"machineHoldReason,omitempty" gorm:"foreignKey:ReasonID;references:ID"`
	ShiftRoster       *ShiftRoster                `json:"shiftRoster,omitempty" gorm:"foreignKey:PausedRosterID;references:ID"`
	HoldRoster        *ShiftRoster                `json:"holdRoster,omitempty" gorm:"foreignKey:PausedRosterID;references:ID"`
	ResumeRoster      *ShiftRoster                `json:"resumeRoster,omitempty" gorm:"foreignKey:ResumeRosterID;references:ID"`
	Machine           *Machine                    `json:"machine,omitempty" gorm:"foreignKey:MachineID"`
	OperationProcess  *OperationProcess           `json:"operationProcess,omitempty" gorm:"foreignKey:JobOperationID;references:ID"`
	Attributes        []OperationProcessAttribute `json:"attributes,omitempty" gorm:"foreignKey:OperationLogID"`
}

func (OperationLog) TableName() string {
	return "tbl_operation_log"
}
