package schema

type ProcessAttribute struct {
	ID            int    `json:"id" gorm:"primaryKey;column:id"`
	ProcessID     int    `json:"processId" gorm:"column:process_id;not null;index"`
	AttributeName string `json:"attributeName" gorm:"type:varchar(255);column:attribute_name;not null"`
	AttributeType string `json:"attributeType" gorm:"type:enum('Output','Waste','Other');default:'Output';not null;column:attribute_type"`
	IsRequired    int    `json:"isRequired" gorm:"type:tinyint(1);default:0;not null;column:is_required"`
	UnitID        *int   `json:"unitId,omitempty" gorm:"column:unit_id;index"`
	// Relationships
	Process *Process `json:"process,omitempty" gorm:"foreignKey:ProcessID"`
	Unit    *Unit    `json:"unit,omitempty" gorm:"foreignKey:UnitID"`
}

func (ProcessAttribute) TableName() string {
	return "tbl_process_attributes"
}
