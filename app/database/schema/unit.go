package schema

import (
	"time"

	"gorm.io/gorm"
)

type Unit struct {
	ID         int            `json:"id" gorm:"primaryKey;column:id"`
	UnitCode   string         `json:"unitCode" gorm:"type:varchar(50);column:unit_code"`
	UnitName   string         `json:"unitName" gorm:"type:varchar(255);column:unit_name"`
	UnitSymbol string         `json:"unitSymbol" gorm:"type:varchar(255);column:unit_symbol"`
	UnitStatus string         `json:"unitStatus" gorm:"type:varchar(255);column:unit_status"`
	CreatedBy  *string        `json:"createdBy" gorm:"type:varchar(255);column:created_by"`
	UpdatedBy  *string        `json:"updatedBy" gorm:"type:varchar(255);column:updated_by"`
	DeletedBy  *string        `json:"deletedBy" gorm:"type:varchar(255);column:deleted_by"`
	CreatedAt  time.Time      `json:"createdAt" gorm:"autoCreateTime;column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"index;column:deleted_at"`

	// Relationships
	Machines          []Machine          `json:"machines,omitempty" gorm:"foreignKey:UnitOfMachineID"`
	ProcessAttributes []ProcessAttribute `json:"processAttributes,omitempty" gorm:"foreignKey:UnitID"`
}

func (Unit) TableName() string {
	return "tbl_unit"
}
