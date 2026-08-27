package schema

import "time"

type MakeReadyLog struct {
	ID                  int        `json:"id" gorm:"primaryKey;column:id"`
	MakeReadyID         int        `json:"makeReadyId" gorm:"column:make_ready_id;not null;index"`
	HoldTime            time.Time  `json:"holdTime" gorm:"column:hold_time;type:datetime(3);default:CURRENT_TIMESTAMP(3)"`
	ResumeTime          *time.Time `json:"resumeTime,omitempty" gorm:"column:resume_time;type:datetime(3)"`
	MachineHoldReasonID *int       `json:"holdReasonId,omitempty" gorm:"column:hold_reason_id;index"`
	HoldRosterID        *int       `json:"holdRosterId,omitempty" gorm:"column:hold_roster_id;index"`
	ResumeRosterID      *int       `json:"resumeRosterId,omitempty" gorm:"column:resume_roster_id"`

	// Relationships
	MakeReady    *MakeReady   `json:"makeReady,omitempty" gorm:"foreignKey:MakeReadyID"`
	HoldRoster   *ShiftRoster `json:"holdRoster,omitempty" gorm:"foreignKey:HoldRosterID;references:ID"`
	ResumeRoster *ShiftRoster `json:"resumeRoster,omitempty" gorm:"foreignKey:ResumeRosterID;references:ID"`
}

func (MakeReadyLog) TableName() string {
	return "tbl_make_ready_log"
}
