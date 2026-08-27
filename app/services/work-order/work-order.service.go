package workorder

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/thimira/production-tracer/app/utils"
	"github.com/thimira/production-tracer/internal/db"
)

// salesOrderPayload is the inbound Frappe sales-order request ({"data": {...}}).
type salesOrderPayload struct {
	Data salesOrder `json:"data"`
}

type salesOrder struct {
	Name         string           `json:"name"`
	CustomerName string           `json:"customer_name"`
	PONo         string           `json:"po_no"`
	Title        string           `json:"title"`
	TotalQty     float64          `json:"total_qty"`
	DeliveryDate string           `json:"delivery_date"`
	Items        []salesOrderItem `json:"items"`
}

type salesOrderItem struct {
	Name                          string  `json:"name"`
	ItemCode                      string  `json:"item_code"`
	ItemName                      string  `json:"item_name"`
	CustomSalesOrderItemReference string  `json:"custom_sales_order_item_reference"`
	Description                   string  `json:"description"`
	ItemGroup                     string  `json:"item_group"`
	Qty                           float64 `json:"qty"`
	StockUOM                      string  `json:"stock_uom"`
	UOM                           string  `json:"uom"`
	DeliveryDate                  string  `json:"delivery_date"`
}

// saveItem is the per-item shape passed to work_order_save as a JSON array.
type saveItem struct {
	ItemReference string  `json:"itemReference"`
	ItemCode      string  `json:"itemCode"`
	ItemName      string  `json:"itemName"`
	UOM           string  `json:"uom"`
	StockUOM      string  `json:"stockUom"`
	ItemGroup     string  `json:"itemGroup"`
	DeliveryDate  string  `json:"deliveryDate"`
	Description   string  `json:"description"`
	Qty           float64 `json:"qty"`
}

// saveResult is the row returned by work_order_save.
type saveResult struct {
	ID                 int64  `json:"id" gorm:"column:id"`
	WorkOrderReference string `json:"workOrderReference" gorm:"column:work_order_reference"`
	ItemCount          int    `json:"itemCount" gorm:"column:item_count"`
}

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

// Save imports/creates a work order and its items from a Frappe sales-order payload.
func Save(c *gin.Context) {
	var payload salesOrderPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "invalid request body: " + err.Error()})
		return
	}

	so := payload.Data
	if strings.TrimSpace(so.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "work order reference (data.name) is required"})
		return
	}

	items := make([]saveItem, 0, len(so.Items))
	for _, it := range so.Items {
		items = append(items, saveItem{
			ItemReference: firstNonEmpty(it.CustomSalesOrderItemReference, it.Name),
			ItemCode:      it.ItemCode,
			ItemName:      it.ItemName,
			UOM:           it.UOM,
			StockUOM:      it.StockUOM,
			ItemGroup:     it.ItemGroup,
			DeliveryDate:  strings.TrimSpace(it.DeliveryDate),
			Description:   stripHTML(it.Description),
			Qty:           it.Qty,
		})
	}
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "failed to encode items: " + err.Error()})
		return
	}

	actor := c.GetHeader("X-User")
	if actor == "" {
		actor = "system"
	}

	result, err := db.CallProcedure[saveResult](
		"work_order_save",
		so.Name,
		nullable(so.CustomerName),
		nullable(so.PONo),
		nullable(so.Title),
		int64(so.TotalQty),
		parseDate(so.DeliveryDate),
		actor,
		string(itemsJSON),
	)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error_message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetWorkOrderById returns a work order with its items.
func GetWorkOrderById(id int) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("work_order_get", id)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "items")
	return result, nil
}

// FindWorkOrders lists work orders by filters with cursor pagination.
func FindWorkOrders(reference, customer, status string, cursor int64, limit int) ([]map[string]interface{}, error) {
	if cursor == 0 {
		cursor = -1
	}
	if limit <= 0 {
		limit = 20
	}
	return db.CallProcedure[[]map[string]interface{}]("work_order_find",
		reference, customer, status, cursor, limit,
	)
}

// UpdateWorkOrder patches a work order header (0 id when not found).
func UpdateWorkOrder(id int, updates map[string]interface{}, updatedBy string) (int, error) {
	existing, err := GetWorkOrderById(id)
	if err != nil {
		return 0, err
	}
	if len(existing) == 0 {
		return 0, nil
	}

	customerName := utils.GetPatchValue(updates, existing, "customerName")
	poNo := utils.GetPatchValue(updates, existing, "poNo")
	title := utils.GetPatchValue(updates, existing, "title")
	quantity := utils.GetPatchValue(updates, existing, "quantity")
	deliveryDate := utils.GetPatchValue(updates, existing, "deliveryDate")
	status := utils.GetPatchValue(updates, existing, "status")

	result, err := db.CallProcedure[map[string]interface{}]("work_order_update",
		id, customerName, poNo, title, quantity, deliveryDate, status, updatedBy,
	)
	if err != nil {
		return 0, err
	}
	return utils.ToInt(result["workOrderId"])
}

func stripHTML(s string) string {
	return strings.TrimSpace(htmlTagRe.ReplaceAllString(s, ""))
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

// nullable returns nil for an empty string so the column is stored as NULL.
func nullable(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

// parseDate parses common Frappe date formats; nil (SQL NULL) if empty/unparseable.
func parseDate(s string) any {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	for _, layout := range []string{"2006-01-02", "2006-01-02 15:04:05", time.RFC3339} {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return nil
}
