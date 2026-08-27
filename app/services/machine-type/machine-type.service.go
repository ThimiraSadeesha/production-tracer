package machinetype

import (
	"encoding/json"

	"github.com/thimira/production-tracer/app/utils"
	"github.com/thimira/production-tracer/internal/db"
)

// GetMachineTypeById returns a machine type with its nested machines.
func GetMachineTypeById(id int) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("machine_type_get", id)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "machines")
	return result, nil
}

// CreateMachineTypes bulk-inserts machine types (and nested machines) from a JSON array (or a single object).
func CreateMachineTypes(payload json.RawMessage, createdBy string) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("machine_type_save", string(payload), createdBy)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "insertedIds")
	return result, nil
}
