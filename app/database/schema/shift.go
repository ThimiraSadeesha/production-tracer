package schema

import (
	"time"

	"gorm.io/gorm"
)

type Shift struct {
	ID           int            `json:"id" gorm:"primaryKey;column:id"`
	Shift        string         `json:"shift" gorm:"type:varchar(255);column:shift"`
	StartDate    string         `json:"startDate" gorm:"type:varchar(10);column:start_date"`
	StartTime    string         `json:"startTime" gorm:"type:varchar(8);column:start_time"` // keep string for API
	EndDate      string         `json:"endDate" gorm:"type:varchar(10);column:end_date"`
	EndTime      string         `json:"endTime" gorm:"type:varchar(8);column:end_time"`
	Status       string         `json:"status" gorm:"type:varchar(100);column:status"`
	DepartmentID int            `json:"departmentId" gorm:"column:department_id"`
	CreatedBy    *string        `json:"createdBy" gorm:"type:varchar(255);column:created_by"`
	UpdatedBy    *string        `json:"updatedBy" gorm:"type:varchar(255);column:updated_by"`
	DeletedBy    *string        `json:"deletedBy" gorm:"type:varchar(255);column:deleted_by"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"autoCreateTime;column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt" gorm:"index;column:deleted_at"`

	// Relationships
	Department   *Department   `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	ShiftRosters []ShiftRoster `json:"shiftRosters,omitempty" gorm:"foreignKey:ShiftID"`
}

func (Shift) TableName() string {
	return "tbl_shift"
}
