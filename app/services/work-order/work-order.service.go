package workorder

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/thimira/production-tracer/app/utils"
	"github.com/thimira/production-tracer/internal/db"
)

// Save bulk-inserts work orders from Frappe sales-order JSON.
// work_order_save expects a JSON array of objects with snake_case fields
// (name, customer_name, po_no, title, total_qty, delivery_date, items).
// Accepted request shapes: a JSON array, a single sales-order object,
// or a Frappe envelope ({"data": {...}} / {"data": [...]}).
func Save(payload json.RawMessage, createdBy string) (map[string]interface{}, error) {
	normalized, err := normalizeSavePayload(payload)
	if err != nil {
		return nil, err
	}
	result, err := db.CallProcedure[map[string]interface{}]("work_order_save", string(normalized), createdBy)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "insertedIds")
	return result, nil
}

func normalizeSavePayload(raw json.RawMessage) (json.RawMessage, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("work order payload cannot be empty")
	}

	// Recover `{ [ ... ] }` — invalid JSON some clients send around an array.
	if !json.Valid(trimmed) && trimmed[0] == '{' && trimmed[len(trimmed)-1] == '}' {
		inner := bytes.TrimSpace(trimmed[1 : len(trimmed)-1])
		if json.Valid(inner) {
			trimmed = inner
		}
	}

	if !json.Valid(trimmed) {
		return nil, fmt.Errorf("invalid request body")
	}

	// Unwrap Frappe envelope { "data": [...] } or { "data": { ... } }.
	if trimmed[0] == '{' {
		var env struct {
			Data json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(trimmed, &env); err == nil {
			if data := bytes.TrimSpace(env.Data); len(data) > 0 && string(data) != "null" {
				trimmed = data
			}
		}
	}

	if len(trimmed) == 0 {
		return nil, fmt.Errorf("work order payload cannot be empty")
	}

	// Wrap a single sales-order object into the array the SP iterates.
	if trimmed[0] == '{' {
		wrapped, err := json.Marshal([]json.RawMessage{json.RawMessage(trimmed)})
		if err != nil {
			return nil, err
		}
		return wrapped, nil
	}

	if trimmed[0] != '[' {
		return nil, fmt.Errorf("work order payload must be a JSON array")
	}
	return trimmed, nil
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
