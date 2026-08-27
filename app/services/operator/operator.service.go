package operator

import (
	"github.com/thimira/production-tracer/app/database/schema"
	"github.com/thimira/production-tracer/app/utils"
	"github.com/thimira/production-tracer/internal/db"
)

// GetOperatorById returns a single operator by id.
func GetOperatorById(id int) (map[string]interface{}, error) {
	return db.CallProcedure[map[string]interface{}]("operator_get", id)
}

// FindOperators lists operators by filters with cursor pagination.
func FindOperators(empNo, name, section, status string, departmentID, cursor int64, limit int) ([]map[string]interface{}, error) {
	if cursor == 0 {
		cursor = -1
	}
	if limit <= 0 {
		limit = 20
	}
	return db.CallProcedure[[]map[string]interface{}]("operator_find",
		empNo, name, section, status, departmentID, cursor, limit,
	)
}

// GetOperatorsByDepartment returns a department with its operators.
func GetOperatorsByDepartment(departmentID int64) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("operator_get_by_department", departmentID)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "operators")
	return result, nil
}

// FindOperatorsByShiftDept returns operators for a department + shift.
func FindOperatorsByShiftDept(departmentID, shiftID int64) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("find_operators_by_shift_dept", departmentID, shiftID)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "operators")
	return result, nil
}

// CreateOperator inserts an operator and returns its id.
func CreateOperator(op schema.Operator, createdBy string) (int, error) {
	result, err := db.CallProcedure[map[string]interface{}]("operator_save",
		op.EmpNo, op.Name, op.Section, op.Status, op.DepartmentID, createdBy,
	)
	if err != nil {
		return 0, err
	}
	return utils.ToInt(result["operatorId"])
}

// UpdateOperator patches an operator (0 id when not found).
func UpdateOperator(id int, updates map[string]interface{}, updatedBy string) (int, error) {
	existing, err := GetOperatorById(id)
	if err != nil {
		return 0, err
	}
	if len(existing) == 0 {
		return 0, nil
	}

	empNo := utils.GetPatchValue(updates, existing, "empNo")
	name := utils.GetPatchValue(updates, existing, "name", "operatorName")
	section := utils.GetPatchValue(updates, existing, "section")
	status := utils.GetPatchValue(updates, existing, "status", "operatorStatus")
	departmentID := utils.GetPatchValue(updates, existing, "departmentId")

	result, err := db.CallProcedure[map[string]interface{}]("operator_update",
		id, empNo, name, section, status, departmentID, updatedBy,
	)
	if err != nil {
		return 0, err
	}
	return utils.ToInt(result["operatorId"])
}
