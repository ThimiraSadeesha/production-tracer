package department

import (
	"github.com/thimira/production-tracer/app/database/schema"
	"github.com/thimira/production-tracer/app/utils"
	"github.com/thimira/production-tracer/internal/db"
)

// GetDepartmentById returns a department with its processes.
func GetDepartmentById(id int) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("department_get", id)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "processes")
	return result, nil
}

// GetDepartments lists all departments.
func GetDepartments() ([]map[string]interface{}, error) {
	return db.CallProcedure[[]map[string]interface{}]("department_get_all")
}

// FindDepartments lists departments by filters with cursor pagination.
func FindDepartments(departmentCode, departmentName, departmentStatus string, cursor int64, limit int) ([]map[string]interface{}, error) {
	if cursor == 0 {
		cursor = -1
	}
	if limit <= 0 {
		limit = 20
	}
	return db.CallProcedure[[]map[string]interface{}]("department_find",
		departmentCode,
		departmentName,
		departmentStatus,
		cursor,
		limit,
	)
}

// CreateDepartment inserts a department and returns its id.
func CreateDepartment(department schema.Department, createdBy string) (int, error) {
	result, err := db.CallProcedure[map[string]interface{}]("department_save",
		department.DepartmentCode,
		department.DepartmentName,
		createdBy,
	)
	if err != nil {
		return 0, err
	}
	return utils.ToInt(result["departmentId"])
}

// UpdateDepartment patches a department (0 id when not found).
func UpdateDepartment(id int, updates map[string]interface{}, updatedBy string) (int, error) {
	existing, err := GetDepartmentById(id)
	if err != nil {
		return 0, err
	}
	if existing == nil {
		return 0, nil
	}

	depCode := utils.GetPatchValue(updates, existing, "departmentCode")
	depName := utils.GetPatchValue(updates, existing, "departmentName")
	depStatus := utils.GetPatchValue(updates, existing, "departmentStatus")

	result, err := db.CallProcedure[map[string]interface{}]("department_update",
		id,
		depCode,
		depName,
		depStatus,
		updatedBy,
	)
	if err != nil {
		return 0, err
	}
	return utils.ToInt(result["departmentId"])
}
