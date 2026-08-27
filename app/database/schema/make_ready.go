package schema

import "time"

type MakeReady struct {
	ID                 int        `json:"id" gorm:"primaryKey;column:id"`
	OperationProcessID int        `json:"operationProcessId" gorm:"column:operation_process_id;not null;index"`
	MachineID          *int       `json:"machineId,omitempty" gorm:"column:machine_id;index"`
	StartedAt          *time.Time `json:"startedAt,omitempty" gorm:"column:started_at;type:datetime(3)"`
	EndedAt            *time.Time `json:"endedAt,omitempty" gorm:"column:ended_at;type:datetime(3)"`
	StartedRosterID    *int       `json:"startedRosterId,omitempty" gorm:"column:started_roster_id;index"`
	EndedRosterID      *int       `json:"endedRosterId,omitempty" gorm:"column:ended_roster_id"`
	Status             string     `json:"status" gorm:"type:varchar(50);column:status;default:InProgress"`
	CreatedAt          time.Time  `json:"createdAt" gorm:"autoCreateTime;column:created_at;type:datetime(3);default:CURRENT_TIMESTAMP(3)"`
	UpdatedAt          time.Time  `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at;type:datetime(3);default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)"`
	Reason             *string    `json:"reason" gorm:"type:varchar(255);column:reason"`

	Machine    *Machine              `json:"machine,omitempty" gorm:"foreignKey:MachineID;references:ID"`
	Logs       []MakeReadyLog        `json:"logs,omitempty" gorm:"foreignKey:MakeReadyID"`
	Attributes []MakeReadyAttributes `json:"attributes,omitempty" gorm:"foreignKey:MakeReadyID"`
}

func (MakeReady) TableName() string {
	return "tbl_make_ready"
}

type MakeReadyLogStatus string

const (
	MakeReadyLogStatusInProgress MakeReadyLogStatus = "Processing"
	MakeReadyLogStatusOnHold     MakeReadyLogStatus = "Paused"
	MakeReadyLogStatusCompleted  MakeReadyLogStatus = "Completed"
)
