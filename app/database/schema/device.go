package schema

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	ID         int            `json:"id" gorm:"primaryKey;column:id"`
	AndroidId  string         `json:"androidId" gorm:"type:varchar(150);column:android_id"`
	AppVersion string         `json:"appVersion" gorm:"type:varchar(50);column:app_version"`
	DeviceName string         `json:"deviceName" gorm:"type:varchar(255);column:device_name"`
	LastSeenAt string         `json:"lastSeenAt" gorm:"column:last_seen_at"`
	Status     string         `json:"status" gorm:"type:varchar(50);column:status"`
	DeviceType string         `json:"deviceType" gorm:"type:enum('Terminal','TV','Manual');default:'Terminal';not null;column:device_type"`
	CreatedBy  *string        `json:"createdBy" gorm:"type:varchar(255);column:created_by"`
	UpdatedBy  *string        `json:"updatedBy" gorm:"type:varchar(255);column:updated_by"`
	DeletedBy  *string        `json:"deletedBy" gorm:"type:varchar(255);column:deleted_by"`
	CreatedAt  time.Time      `json:"createdAt" gorm:"autoCreateTime;column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"index;column:deleted_at"`

	// Optional relationships
	MachineID *int     `json:"machineId,omitempty" gorm:"column:machine_id"` // nullable
	ProcessID *int     `json:"processId,omitempty" gorm:"column:process_id"` // nullable
	Machine   *Machine `json:"machine,omitempty" gorm:"foreignKey:MachineID"`
	Process   *Process `json:"process,omitempty" gorm:"foreignKey:ProcessID"`
}

func (Device) TableName() string {
	return "tbl_device"
}
