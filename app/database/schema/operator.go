package schema

import (
	"time"

	"gorm.io/gorm"
)

type Operator struct {
	ID           int            `json:"id" gorm:"primaryKey;column:id"`
	Status       string         `json:"status" gorm:"type:varchar(100);column:status"`
	EmpNo        string         `json:"empNo" gorm:"type:varchar(255);column:emp_no"`
	Name         string         `json:"name" gorm:"type:varchar(255);column:name"`
	Section      string         `json:"section" gorm:"type:varchar(255);column:section"`
	DepartmentID int            `json:"departmentId" gorm:"column:department_id"`
	CreatedBy    *string        `json:"createdBy" gorm:"type:varchar(255);column:created_by"`
	UpdatedBy    *string        `json:"updatedBy" gorm:"type:varchar(255);column:updated_by"`
	DeletedBy    *string        `json:"deletedBy" gorm:"type:varchar(255);column:deleted_by"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"autoCreateTime;column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt" gorm:"index;column:deleted_at"`

	// Relationships
	Department   *Department   `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	ShiftRosters []ShiftRoster `json:"shiftRosters,omitempty" gorm:"foreignKey:OperatorID"`
}

func (Operator) TableName() string {
	return "tbl_operator"
}
