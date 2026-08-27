package schema

import "time"

type MakeReadyAttributes struct {
	ID          int       `json:"id" gorm:"primaryKey;column:id"`
	MakeReadyID int       `json:"makeReadyId" gorm:"column:make_ready_id;not null;index"`
	AttributeID int       `json:"attributeId" gorm:"column:attribute_id;not null;index"`
	Value       int       `json:"value" gorm:"column:value;default:0;not null"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`

	// Relationships
	MakeReady *MakeReady        `json:"makeReady,omitempty" gorm:"foreignKey:MakeReadyID"`
	Attribute *ProcessAttribute `json:"attribute,omitempty" gorm:"foreignKey:AttributeID"`
}

func (MakeReadyAttributes) TableName() string {
	return "tbl_make_ready_attributes"
}
