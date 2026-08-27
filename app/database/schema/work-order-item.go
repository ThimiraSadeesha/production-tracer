package schema

import "time"

// WorkOrderItem maps to tbl_work_order_item — one line of a work order,
// imported from a Frappe sales order item (see dtos.SalesOrderItem).
type WorkOrderItem struct {
	ID            int        `json:"id" gorm:"primaryKey;column:id"`
	WorkOrderID   int        `json:"workOrderId" gorm:"index;not null;column:work_order_id"`
	ItemReference string     `json:"itemReference" gorm:"type:varchar(255);column:item_reference"`
	ItemCode      string     `json:"itemCode" gorm:"type:varchar(255);column:item_code"`
	ItemName      string     `json:"itemName" gorm:"type:varchar(255);column:item_name"`
	UOM           string     `json:"uom" gorm:"type:varchar(50);column:uom"`
	DeliveryDate  *time.Time `json:"deliveryDate" gorm:"column:delivery_date"`
	Description   string     `json:"description" gorm:"type:varchar(255);column:remarks"`
	ItemGroup     string     `json:"itemGroup" gorm:"type:varchar(255);column:item_group"`
	Quantity      int64      `json:"qty" gorm:"default:0;column:quantity"`
	StockUOM      string     `json:"stockUom" gorm:"type:varchar(50);column:stock_uom"`
	ItemStatus    string     `json:"itemStatus" gorm:"type:varchar(100);default:Pending;column:item_status"`

	Audit

	// Relationships
	WorkOrder *WorkOrder `json:"workOrder,omitempty" gorm:"foreignKey:WorkOrderID"`
}

func (WorkOrderItem) TableName() string {
	return "tbl_work_order_item"
}
