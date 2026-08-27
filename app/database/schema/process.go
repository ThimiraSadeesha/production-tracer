package schema

import (
	"time"

	"gorm.io/gorm"
)

type Process struct {
	ID           int            `json:"id" gorm:"primaryKey;column:id"`
	Code         string         `json:"code" gorm:"type:varchar(100);column:code"`
	Name         string         `json:"name" gorm:"type:varchar(255);column:name"`
	Status       string         `json:"status" gorm:"type:varchar(100);column:status"`
	Sequence     int            `json:"sequence" gorm:"column:sequence;default:0"`
	IsRequired   int            `json:"isRequired" gorm:"type:tinyint(1);default:0;not null;column:is_required"`
	Target       *float64       `json:"target" gorm:"column:target;default:0"`
	DepartmentID int            `json:"departmentId" gorm:"column:department_id"`
	CreatedBy    *string        `json:"createdBy" gorm:"type:varchar(255);column:created_by"`
	UpdatedBy    *string        `json:"updatedBy" gorm:"type:varchar(255);column:updated_by"`
	DeletedBy    *string        `json:"deletedBy" gorm:"type:varchar(255);column:deleted_by"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"autoCreateTime;column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt" gorm:"index;column:deleted_at"`

	// Relationships
	Department        *Department        `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	MachineTypes      []MachineType      `json:"machineTypes,omitempty" gorm:"foreignKey:ProcessID"`
	ProcessAttributes []ProcessAttribute `json:"processAttributes,omitempty" gorm:"foreignKey:ProcessID"`
}

func (Process) TableName() string {
	return "tbl_process"
}
