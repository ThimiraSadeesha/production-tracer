package device

import (
	"encoding/json"

	"github.com/thimira/production-tracer/app/utils"
	"github.com/thimira/production-tracer/internal/db"
)

// GetDeviceById returns a device with its machines and processes.
func GetDeviceById(id int) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("device_get_by_id", id)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "machines", "processes")
	return result, nil
}

// FindDevices lists devices by filters with cursor pagination.
func FindDevices(androidID string, processID int64, status string, machineID, cursor int64, limit int) ([]map[string]interface{}, error) {
	if cursor == 0 {
		cursor = -1
	}
	if limit <= 0 {
		limit = 20
	}
	return db.CallProcedure[[]map[string]interface{}]("device_find",
		androidID, processID, status, machineID, cursor, limit,
	)
}

// IdentifyDevice heartbeats a device by android id and returns its resolved context.
func IdentifyDevice(androidID string) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("device_identify", androidID)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "machines")
	return result, nil
}

// CreateDevices bulk-inserts devices from a JSON array; returns saved count and ids.
func CreateDevices(payload json.RawMessage, createdBy string) (map[string]interface{}, error) {
	result, err := db.CallProcedure[map[string]interface{}]("device_save", string(payload), createdBy)
	if err != nil {
		return nil, err
	}
	utils.ProcessJsonData[[]interface{}](result, "insertedIds")
	return result, nil
}

// UpdateDevice patches a device (0 id when not found).
func UpdateDevice(id int, updates map[string]interface{}, updatedBy string) (int, error) {
	existing, err := GetDeviceById(id)
	if err != nil {
		return 0, err
	}
	if len(existing) == 0 {
		return 0, nil
	}

	appVersion := utils.GetPatchValue(updates, existing, "appVersion")
	lastSeenAt := utils.GetPatchValue(updates, existing, "lastSeenAt")
	status := utils.GetPatchValue(updates, existing, "status", "deviceStatus")
	deviceType := utils.GetPatchValue(updates, existing, "deviceType")
	machineID := utils.GetPatchValue(updates, existing, "machineId")
	processID := utils.GetPatchValue(updates, existing, "processId")
	androidID := utils.GetPatchValue(updates, existing, "androidId")

	result, err := db.CallProcedure[map[string]interface{}]("device_update",
		id, appVersion, lastSeenAt, status, deviceType, machineID, processID, androidID, updatedBy,
	)
	if err != nil {
		return 0, err
	}
	return utils.ToInt(result["deviceId"])
}
