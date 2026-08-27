package schema

import (
	"time"

	"gorm.io/gorm"
)

type MachineType struct {
	ID        int            `json:"id" gorm:"primaryKey;column:id"`
	Code      string         `json:"code" gorm:"type:varchar(100);column:code"`
	Name      string         `json:"name" gorm:"type:varchar(255);column:name"`
	Status    string         `json:"status" gorm:"type:varchar(100);column:machine_type_status;default:'Active'"`
	ProcessID int            `json:"processId" gorm:"column:process_id"`
	CreatedBy *string        `json:"createdBy" gorm:"type:varchar(255);column:created_by"`
	UpdatedBy *string        `json:"updatedBy" gorm:"type:varchar(255);column:updated_by"`
	DeletedBy *string        `json:"deletedBy" gorm:"type:varchar(255);column:deleted_by"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime;column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index;column:deleted_at"`
	Process   *Process       `json:"process,omitempty" gorm:"foreignKey:ProcessID"`
	Machines  []Machine      `json:"machines,omitempty" gorm:"foreignKey:MachineTypeID"`
}

func (MachineType) TableName() string {
	return "tbl_machine_type"
}
