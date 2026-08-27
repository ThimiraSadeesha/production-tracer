package unit

import (
	"github.com/thimira/production-tracer/app/database/schema"
	"github.com/thimira/production-tracer/app/utils"
	"github.com/thimira/production-tracer/internal/db"
)

// GetUnits lists all units.
func GetUnits() ([]map[string]interface{}, error) {
	return db.CallProcedure[[]map[string]interface{}]("unit_getAll")
}

// GetUnitById returns a single unit by id.
func GetUnitById(id int) (map[string]interface{}, error) {
	return db.CallProcedure[map[string]interface{}]("unit_get", id)
}

// CreateUnit inserts a unit and returns its id.
func CreateUnit(unit schema.Unit, createdBy string) (int, error) {
	result, err := db.CallProcedure[map[string]interface{}]("unit_save",
		unit.UnitCode,
		unit.UnitName,
		unit.UnitSymbol,
		createdBy,
	)
	if err != nil {
		return 0, err
	}
	return utils.ToInt(result["unitId"])
}

// UpdateUnit patches an existing unit (nil id when not found).
func UpdateUnit(id int, updates map[string]interface{}, updatedBy string) (int, error) {
	existing, err := GetUnitById(id)
	if err != nil {
		return 0, err
	}
	if len(existing) == 0 {
		return 0, nil
	}

	code := utils.GetPatchValue(updates, existing, "unitCode")
	name := utils.GetPatchValue(updates, existing, "unitName")
	symbol := utils.GetPatchValue(updates, existing, "unitSymbol")
	status := utils.GetPatchValue(updates, existing, "unitStatus")

	result, err := db.CallProcedure[map[string]interface{}]("unit_update",
		id, code, name, symbol, status, updatedBy,
	)
	if err != nil {
		return 0, err
	}
	return utils.ToInt(result["unitId"])
}
