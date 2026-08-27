package shift

import (
	"encoding/json"

	"github.com/thimira/production-tracer/app/utils"
	"github.com/thimira/production-tracer/internal/db"
)

// GetShifts lists all shifts.
func GetShifts() ([]map[string]interface{}, error) {
	return db.CallProcedure[[]map[string]interface{}]("shift_getAll")
}

// GetShiftById returns a single shift by id.
func GetShiftById(id int) (map[string]interface{}, error) {
	return db.CallProcedure[map[string]interface{}]("shift_get", id)
}

// GetShiftsByDepartment returns a department with its shifts.
func GetShiftsByDepartment(departmentID int) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("shift_get_by_department", departmentID)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "shifts")
	return result, nil
}

// GetShiftOutputHistory returns one shift's output history in a date range.
func GetShiftOutputHistory(shiftID int64, from, to string) ([]map[string]interface{}, error) {
	return db.CallProcedure[[]map[string]interface{}]("shift_output_history_get",
		shiftID, utils.EmptyToNil(from), utils.EmptyToNil(to),
	)
}

// FindShiftOutputHistory returns output history filtered by department/shift/date.
func FindShiftOutputHistory(departmentID, shiftID int64, from, to string) ([]map[string]interface{}, error) {
	return db.CallProcedure[[]map[string]interface{}]("shift_output_history_find",
		departmentID, shiftID, utils.EmptyToNil(from), utils.EmptyToNil(to),
	)
}

// CreateShifts bulk-inserts shifts from a JSON array (or a single object).
func CreateShifts(payload json.RawMessage, createdBy string) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("shift_save", string(payload), createdBy)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "insertedIds")
	return result, nil
}

// UpdateShift patches a shift (0 id when not found; SP returns no row on success).
func UpdateShift(id int, updates map[string]interface{}, updatedBy string) (int, error) {
	existing, err := GetShiftById(id)
	if err != nil {
		return 0, err
	}
	if len(existing) == 0 {
		return 0, nil
	}

	name := utils.GetPatchValue(updates, existing, "shiftName", "shift")
	startTime := utils.GetPatchValue(updates, existing, "startTime")
	endTime := utils.GetPatchValue(updates, existing, "endTime")
	departmentID := utils.GetPatchValue(updates, existing, "departmentId")
	status := utils.GetPatchValue(updates, existing, "status", "shiftStatus")

	if _, err := db.CallProcedure[map[string]interface{}]("shift_update",
		id, name, startTime, endTime, departmentID, status, updatedBy,
	); err != nil {
		return 0, err
	}
	return id, nil
}
