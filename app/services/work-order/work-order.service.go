package workorder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/thimira/production-tracer/app/utils"
	"github.com/thimira/production-tracer/internal/db"
)

// salesOrder is the inbound Frappe sales-order shape (snake_case).
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

// nativeWorkOrder is the per-order shape passed to work_order_save.
type nativeWorkOrder struct {
	WorkOrderReference string     `json:"workOrderReference"`
	CustomerName       string     `json:"customerName"`
	PONo               string     `json:"poNo"`
	Title              string     `json:"title"`
	TotalQty           float64    `json:"totalQty"`
	DeliveryDate       string     `json:"deliveryDate"`
	Items              []saveItem `json:"items"`
}

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

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

// Save imports/creates work orders and their items from a JSON array, a single
// object, or a Frappe envelope ({"data": {...}} / {"data": [...]}).
func Save(c *gin.Context) {
	var raw json.RawMessage
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "invalid request body: " + err.Error()})
		return
	}

	payload, err := normalizeWorkOrderSavePayload(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
		return
	}

	actor := c.GetHeader("X-User")
	if actor == "" {
		actor = "system"
	}

	result, err := db.CallProcedure[map[string]interface{}]("work_order_save", string(payload), actor)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error_message": err.Error()})
		return
	}
	utils.ProcessJsonData[[]interface{}](result, "insertedIds")
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func normalizeWorkOrderSavePayload(raw json.RawMessage) (json.RawMessage, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("invalid request body")
	}

	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(trimmed, &env); err == nil && len(bytes.TrimSpace(env.Data)) > 0 {
		orders, err := mapFrappeSalesOrders(env.Data)
		if err != nil {
			return nil, err
		}
		for _, o := range orders {
			if strings.TrimSpace(o.WorkOrderReference) == "" {
				return nil, fmt.Errorf("work order reference (data.name) is required")
			}
		}
		return json.Marshal(orders)
	}

	return trimmed, nil
}

func mapFrappeSalesOrders(data json.RawMessage) ([]nativeWorkOrder, error) {
	trimmed := bytes.TrimSpace(data)
	var sos []salesOrder
	if len(trimmed) > 0 && trimmed[0] == '[' {
		if err := json.Unmarshal(trimmed, &sos); err != nil {
			return nil, fmt.Errorf("invalid request body: %w", err)
		}
	} else {
		var so salesOrder
		if err := json.Unmarshal(trimmed, &so); err != nil {
			return nil, fmt.Errorf("invalid request body: %w", err)
		}
		sos = []salesOrder{so}
	}

	out := make([]nativeWorkOrder, 0, len(sos))
	for _, so := range sos {
		out = append(out, mapSalesOrder(so))
	}
	return out, nil
}

func mapSalesOrder(so salesOrder) nativeWorkOrder {
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
	return nativeWorkOrder{
		WorkOrderReference: so.Name,
		CustomerName:       so.CustomerName,
		PONo:               so.PONo,
		Title:              so.Title,
		TotalQty:           so.TotalQty,
		DeliveryDate:       strings.TrimSpace(so.DeliveryDate),
		Items:              items,
	}
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
