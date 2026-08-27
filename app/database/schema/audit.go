package schema

import (
	"time"

	"gorm.io/gorm"
)

// Audit holds the standard bookkeeping columns shared by every table.
type Audit struct {
	CreatedBy *string        `json:"createdBy" gorm:"type:varchar(255);column:created_by"`
	UpdatedBy *string        `json:"updatedBy" gorm:"type:varchar(255);column:updated_by"`
	DeletedBy *string        `json:"deletedBy" gorm:"type:varchar(255);column:deleted_by"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime;column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index;column:deleted_at"`
}
