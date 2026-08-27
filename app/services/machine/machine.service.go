package machine

import (
	"encoding/json"

	"github.com/thimira/production-tracer/app/utils"
	"github.com/thimira/production-tracer/internal/db"
)

// GetMachineById returns a machine with unit, hold reasons, performance, and current operation.
func GetMachineById(id int) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("machine_get", id)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[map[string]interface{}](result, "unit", "performance", "current_operation")
	utils.ProcessJsonData[[]interface{}](result, "hold_reasons")
	return result, nil
}

// FindMachines lists machines by filters with cursor pagination.
func FindMachines(code, name string, machineTypeID, departmentID, processID, cursor int64, limit int) ([]map[string]interface{}, error) {
	if cursor == 0 {
		cursor = -1
	}
	if limit <= 0 {
		limit = 20
	}
	return db.CallProcedure[[]map[string]interface{}]("machine_find",
		code, name, machineTypeID, departmentID, processID, cursor, limit,
	)
}

// CreateMachines bulk-inserts machines from a JSON array (or a single object).
func CreateMachines(payload json.RawMessage, createdBy string) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("machine_save", string(payload), createdBy)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "insertedIds")
	return result, nil
}

// UpdateMachine patches a machine (0 id when not found).
func UpdateMachine(id int, updates map[string]interface{}, updatedBy string) (int, error) {
	existing, err := GetMachineById(id)
	if err != nil {
		return 0, err
	}
	if len(existing) == 0 {
		return 0, nil
	}

	machineCode := utils.GetPatchValue(updates, existing, "machineCode", "machine_code")
	machineName := utils.GetPatchValue(updates, existing, "machineName", "machine_name")
	description := utils.GetPatchValue(updates, existing, "description")
	capabilities := utils.GetPatchValue(updates, existing, "capabilities")
	status := utils.GetPatchValue(updates, existing, "status", "machine_status")
	hourlyOutput := utils.GetPatchValue(updates, existing, "hourlyOutput", "hourly_output")
	makeReadyTime := utils.GetPatchValue(updates, existing, "makeReadyTime", "make_ready_time")
	machineTypeID := utils.GetPatchValue(updates, existing, "machineTypeId", "machine_type_id")
	unitID := utils.GetPatchValue(updates, existing, "unitOfMachineId", "unit_of_machine_id")

	result, err := db.CallProcedure[map[string]interface{}]("machine_update",
		id, machineCode, machineName, description, capabilities, status,
		hourlyOutput, makeReadyTime, machineTypeID, unitID, updatedBy,
	)
	if err != nil {
		return 0, err
	}
	return utils.ProcedureIntFromMap(result, "updatedMachineId")
}
