package process

import (
	"encoding/json"

	"github.com/thimira/production-tracer/app/utils"
	"github.com/thimira/production-tracer/internal/db"
)

// GetProcesses lists all processes.
func GetProcesses() ([]map[string]interface{}, error) {
	return db.CallProcedure[[]map[string]interface{}]("process_getAll")
}

// GetProcessById returns a process with its machine types and attributes.
func GetProcessById(id int) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("process_get", id)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "machineTypes", "attributes")
	return result, nil
}

// GetProcessesByDepartment returns a department with its processes.
func GetProcessesByDepartment(departmentID int64) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("process_get_by_department", departmentID)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "processes")
	return result, nil
}

// CreateProcesses bulk-inserts processes (and their attributes) from a JSON array.
func CreateProcesses(payload json.RawMessage, createdBy string) (int, error) {
	result, err := db.CallProcedure[map[string]interface{}]("process_save", string(payload), createdBy)
	if err != nil {
		return 0, err
	}
	return utils.ToInt(result["numberOfProcess"])
}

// UpdateProcess patches a process (0 id when not found).
func UpdateProcess(id int, updates map[string]interface{}, updatedBy string) (int, error) {
	existing, err := GetProcessById(id)
	if err != nil {
		return 0, err
	}
	if len(existing) == 0 {
		return 0, nil
	}

	code := utils.GetPatchValue(updates, existing, "code", "processCode")
	name := utils.GetPatchValue(updates, existing, "name", "processName")
	status := utils.GetPatchValue(updates, existing, "status", "processStatus")
	isRequired := utils.GetPatchValue(updates, existing, "isRequired")
	departmentID := utils.GetPatchValue(updates, existing, "departmentId")

	if _, err := db.CallProcedure[map[string]interface{}]("process_update",
		id, code, name, status, isRequired, departmentID, updatedBy,
	); err != nil {
		return 0, err
	}
	return id, nil
}
