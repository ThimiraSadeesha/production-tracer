package schema

import "time"

// WorkOrder maps to tbl_work_order. It is created from an imported Frappe sales
// order (see dtos.SalesOrder) or created manually.
type WorkOrder struct {
	ID                 int        `json:"id" gorm:"primaryKey;column:id"`
	WorkOrderReference string     `json:"workOrderReference" gorm:"type:varchar(255);uniqueIndex;column:work_order_reference"`
	CustomerName       string     `json:"customerName" gorm:"type:varchar(255);column:customer_name"`
	PONo               string     `json:"poNo" gorm:"type:varchar(255);column:po_no"`
	Title              string     `json:"title" gorm:"type:varchar(255);column:title"`
	Quantity           int64      `json:"totalQty" gorm:"default:0;column:quantity"`
	DeliveryDate       *time.Time `json:"deliveryDate" gorm:"column:delivery_date"`
	Status             string     `json:"status" gorm:"type:varchar(100);default:Pending;column:status"`
	NoOfBreakdowns     int        `json:"noOfBreakdowns" gorm:"column:no_of_breakdowns;default:0"`

	Audit

	// Relationships
	Items []WorkOrderItem `json:"items,omitempty" gorm:"foreignKey:WorkOrderID"`
}

func (WorkOrder) TableName() string {
	return "tbl_work_order"
}
