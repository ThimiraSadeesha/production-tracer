package schema

type OperationProcessAttribute struct {
	ID                 int               `json:"id" gorm:"primaryKey;column:id"`
	OperationProcessID *int              `json:"operationProcessId,omitempty" gorm:"column:operation_process_id;index"`
	OperationLogID     *int              `json:"operationLogId,omitempty" gorm:"column:operation_log_id;index"`
	AttributeID        int               `json:"attributeId" gorm:"column:attribute_id;not null;index"`
	Value              int               `json:"value" gorm:"column:value;default:0;not null"`
	OperationProcess   *OperationProcess `json:"operationProcess,omitempty" gorm:"foreignKey:OperationProcessID"`
	OperationLog       *OperationLog     `json:"operationLog,omitempty" gorm:"foreignKey:OperationLogID"`
	Attribute          *ProcessAttribute `json:"attribute,omitempty" gorm:"foreignKey:AttributeID"`
}

func (OperationProcessAttribute) TableName() string {
	return "tbl_operation_process_attributes"
}
