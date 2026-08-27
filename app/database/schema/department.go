package schema

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	ID               int            `json:"id" gorm:"primaryKey;column:id"`
	DepartmentCode   string         `json:"departmentCode" gorm:"type:varchar(50);column:department_code"`
	DepartmentName   string         `json:"departmentName" gorm:"type:varchar(255);column:department_name"`
	DepartmentStatus string         `json:"departmentStatus" gorm:"type:varchar(255);column:department_status"`
	CreatedBy        *string        `json:"createdBy" gorm:"type:varchar(255);column:created_by"`
	UpdatedBy        *string        `json:"updatedBy" gorm:"type:varchar(255);column:updated_by"`
	DeletedBy        *string        `json:"deletedBy" gorm:"type:varchar(255);column:deleted_by"`
	CreatedAt        time.Time      `json:"createdAt" gorm:"autoCreateTime;column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time      `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt        gorm.DeletedAt `json:"deletedAt" gorm:"index;column:deleted_at"`

	// Relationships
	Processes []Process  `json:"processes,omitempty" gorm:"foreignKey:DepartmentID"`
	Operators []Operator `json:"operators,omitempty" gorm:"foreignKey:DepartmentID"`
	Shifts    []Shift    `json:"shifts,omitempty" gorm:"foreignKey:DepartmentID"`
}

func (Department) TableName() string {
	return "tbl_department"
}
