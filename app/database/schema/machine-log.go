package schema

import "time"

type MachineLog struct {
	ID              int        `json:"id" gorm:"primaryKey;column:id"`
	MachineID       *int       `json:"machineId" gorm:"column:machine_id;index:idx_ml_job_end_machine_status,priority:3"`
	MachineStatus   string     `json:"machineStatus" gorm:"type:varchar(100);column:machine_status;index:idx_ml_job_end_machine_status,priority:4"`
	StartedAt       *time.Time `json:"startedAt" gorm:"column:started_at"`
	EndedAt         *time.Time `json:"endedAt" gorm:"column:ended_at;index:idx_ml_job_end_machine_status,priority:2"`
	JobOperationID  *int       `json:"jobOperationId" gorm:"column:job_operation_id;index:idx_ml_job_end_machine_status,priority:1"`
	StartedRosterID *int       `json:"startedRosterId" gorm:"column:started_roaster_id"`
	EndRosterID     *int       `json:"endRosterId" gorm:"column:end_roaster_id"`
	Reason          string     `json:"reason" gorm:"type:longtext;column:reason"`
	Remark          *string    `json:"remark" gorm:"type:varchar(255);column:remark"`

	Machine       *Machine          `json:"machine,omitempty" gorm:"foreignKey:MachineID"`
	Operation     *OperationProcess `json:"operation,omitempty" gorm:"foreignKey:JobOperationID"`
	StartedRoster *ShiftRoster      `json:"startedRoster,omitempty" gorm:"foreignKey:StartedRosterID"`
	EndRoster     *ShiftRoster      `json:"endRoster,omitempty" gorm:"foreignKey:EndRosterID"`
}

func (MachineLog) TableName() string {
	return "tbl_machine_log"
}
