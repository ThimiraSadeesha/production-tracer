package operationprocess

import (
	"encoding/json"

	"github.com/thimira/production-tracer/app/utils"
	"github.com/thimira/production-tracer/internal/db"
)

// Save bulk-creates operation processes from a JSON array (ORDER/ITEMS scope).
func Save(payload json.RawMessage, createdBy string) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("operation_process_save", string(payload), createdBy)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "InsertedIds")
	return result, nil
}

// GetOperationProcessById returns one operation process with its attributes.
func GetOperationProcessById(id int) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("operation_process_get", id)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "attributes")
	return result, nil
}

// FindOperationProcesses lists operation processes by filters with cursor pagination.
func FindOperationProcesses(machineID, processID int64, plannedDate, workOrderID, status string, cursor int64, limit int) ([]map[string]interface{}, error) {
	if cursor == 0 {
		cursor = -1
	}
	if limit <= 0 {
		limit = 20
	}
	return db.CallProcedure[[]map[string]interface{}]("operation_process_find",
		machineID, processID, plannedDate, workOrderID, status, cursor, limit,
	)
}

// UpdateOperationProcess patches an operation process (0 id when not found).
func UpdateOperationProcess(id int, updates map[string]interface{}, updatedBy string) (int, error) {
	existing, err := GetOperationProcessById(id)
	if err != nil {
		return 0, err
	}
	if len(existing) == 0 {
		return 0, nil
	}

	sequence := utils.GetPatchValue(updates, existing, "sequence")
	planDateTime := utils.GetPatchValue(updates, existing, "planDateTime")
	plannedQuantity := utils.GetPatchValue(updates, existing, "plannedQuantity")
	expectedValue := utils.GetPatchValue(updates, existing, "expectedValue")
	estimateTime := utils.GetPatchValue(updates, existing, "estimateTime")
	machineID := utils.GetPatchValue(updates, existing, "machineId")
	processID := utils.GetPatchValue(updates, existing, "processId")
	operationScope := utils.GetPatchValue(updates, existing, "operationScope")
	status := utils.GetPatchValue(updates, existing, "status")
	remark := utils.GetPatchValue(updates, existing, "remark")

	result, err := db.CallProcedure[map[string]interface{}]("operation_process_update",
		id, sequence, planDateTime, plannedQuantity, expectedValue, estimateTime,
		machineID, processID, operationScope, status, remark, updatedBy,
	)
	if err != nil {
		return 0, err
	}
	return utils.ToInt(result["operationProcessId"])
}
