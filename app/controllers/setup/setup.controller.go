// Package setup exposes the shop reference/configuration (master-data) routes.
package setup

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/thimira/production-tracer/app/database/schema"
	departmentSvc "github.com/thimira/production-tracer/app/services/department"
	deviceSvc "github.com/thimira/production-tracer/app/services/device"
	machineSvc "github.com/thimira/production-tracer/app/services/machine"
	machineTypeSvc "github.com/thimira/production-tracer/app/services/machine-type"
	operatorSvc "github.com/thimira/production-tracer/app/services/operator"
	processSvc "github.com/thimira/production-tracer/app/services/process"
	shiftSvc "github.com/thimira/production-tracer/app/services/shift"
	unitSvc "github.com/thimira/production-tracer/app/services/unit"
	"github.com/thimira/production-tracer/internal/helper"
)

// Init registers all /setup master-data routes on the router.
func Init(router *gin.Engine) {
	setup := router.Group("/setup")

	departments := setup.Group("/departments")
	{
		departments.GET("", GetDepartments)
		departments.GET("/find", FindDepartments)
		departments.GET("/:id", GetDepartment)
		departments.POST("", CreateDepartments)
		departments.PATCH("/:id", UpdateDepartment)
	}

	units := setup.Group("/units")
	{
		units.GET("", GetUnits)
		units.GET("/:id", GetUnit)
		units.POST("", CreateUnit)
		units.PATCH("/:id", UpdateUnit)
	}

	processes := setup.Group("/processes")
	{
		processes.GET("", GetProcesses)
		processes.GET("/department/:departmentId", GetProcessesByDepartment)
		processes.GET("/:id", GetProcess)
		processes.POST("", CreateProcesses)
		processes.PATCH("/:id", UpdateProcess)
	}

	operators := setup.Group("/operators")
	{
		operators.GET("/find", FindOperators)
		operators.GET("/by-shift", FindOperatorsByShiftDept)
		operators.GET("/department/:departmentId", GetOperatorsByDepartment)
		operators.GET("/:id", GetOperator)
		operators.POST("", CreateOperators)
		operators.PATCH("/:id", UpdateOperator)
	}

	shifts := setup.Group("/shifts")
	{
		shifts.GET("", GetShifts)
		shifts.GET("/output-history", FindShiftOutputHistory)
		shifts.GET("/department/:departmentId", GetShiftsByDepartment)
		shifts.GET("/:id", GetShift)
		shifts.GET("/:id/output-history", GetShiftOutputHistory)
		shifts.POST("", CreateShifts)
		shifts.PATCH("/:id", UpdateShift)
	}

	devices := setup.Group("/devices")
	{
		devices.GET("/find", FindDevices)
		devices.GET("/identify/:androidId", IdentifyDevice)
		devices.GET("/:id", GetDevice)
		devices.POST("", CreateDevices)
		devices.PATCH("/:id", UpdateDevice)
	}

	machineTypes := setup.Group("/machine-types")
	{
		machineTypes.GET("/:id", GetMachineType)
		machineTypes.POST("", CreateMachineTypes)
	}

	machines := setup.Group("/machines")
	{
		machines.GET("/find", FindMachines)
		machines.GET("/:id", GetMachine)
		machines.POST("", CreateMachines)
		machines.PATCH("/:id", UpdateMachine)
	}
}

// GetDepartments lists all departments.
func GetDepartments(c *gin.Context) {
	departments, err := departmentSvc.GetDepartments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": departments})
}

// GetDepartment returns one department by id.
func GetDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department id"})
		return
	}
	dep, err := departmentSvc.GetDepartmentById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(dep) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": dep})
}

// FindDepartments lists departments by filters with cursor pagination.
func FindDepartments(c *gin.Context) {
	cursor, _ := strconv.ParseInt(c.Query("cursor"), 10, 64)
	limit, _ := strconv.Atoi(c.Query("limit"))
	departments, err := departmentSvc.FindDepartments(c.Query("code"), c.Query("name"), c.Query("status"), cursor, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": departments})
}

// CreateDepartments bulk-creates departments from a JSON array (or a single object).
func CreateDepartments(c *gin.Context) {
	var payload json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	result, err := departmentSvc.CreateDepartments(payload, helper.ResolveActor(c, helper.ActorFromJSON(payload)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Departments created successfully", "data": result})
}

// UpdateDepartment patches a department.
func UpdateDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department id"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	depID, err := departmentSvc.UpdateDepartment(id, updates, helper.ResolveActor(c, helper.ActorFromMap(updates)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if depID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Department updated successfully", "departmentId": depID})
}

// GetUnits lists all units.
func GetUnits(c *gin.Context) {
	units, err := unitSvc.GetUnits()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": units})
}

// GetUnit returns one unit by id.
func GetUnit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid unit id"})
		return
	}
	u, err := unitSvc.GetUnitById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(u) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unit not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": u})
}

// CreateUnit creates a unit.
func CreateUnit(c *gin.Context) {
	var dto schema.Unit
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if dto.UnitCode == "" || dto.UnitName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unit code and name are required"})
		return
	}
	createdBy := ""
	if dto.CreatedBy != nil {
		createdBy = *dto.CreatedBy
	}
	id, err := unitSvc.CreateUnit(dto, helper.ResolveActor(c, createdBy))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Unit created successfully", "unitId": id})
}

// UpdateUnit patches a unit.
func UpdateUnit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid unit id"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	unitID, err := unitSvc.UpdateUnit(id, updates, helper.ResolveActor(c, helper.ActorFromMap(updates)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if unitID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unit not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Unit updated successfully", "unitId": unitID})
}

// GetProcesses lists all processes.
func GetProcesses(c *gin.Context) {
	processes, err := processSvc.GetProcesses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": processes})
}

// GetProcess returns one process by id.
func GetProcess(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid process id"})
		return
	}
	p, err := processSvc.GetProcessById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(p) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Process not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// GetProcessesByDepartment lists a department's processes.
func GetProcessesByDepartment(c *gin.Context) {
	deptID, err := strconv.ParseInt(c.Param("departmentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department id"})
		return
	}
	p, err := processSvc.GetProcessesByDepartment(deptID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// CreateProcesses bulk-creates processes from a JSON array.
func CreateProcesses(c *gin.Context) {
	var payload json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	count, err := processSvc.CreateProcesses(payload, helper.ResolveActor(c, helper.ActorFromJSON(payload)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Processes created successfully", "numberOfProcess": count})
}

// UpdateProcess patches a process.
func UpdateProcess(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid process id"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	processID, err := processSvc.UpdateProcess(id, updates, helper.ResolveActor(c, helper.ActorFromMap(updates)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if processID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Process not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Process updated successfully", "processId": processID})
}

// GetOperator returns one operator by id.
func GetOperator(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operator id"})
		return
	}
	op, err := operatorSvc.GetOperatorById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(op) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Operator not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": op})
}

// FindOperators lists operators by filters with cursor pagination.
func FindOperators(c *gin.Context) {
	operators, err := operatorSvc.FindOperators(
		c.Query("empNo"), c.Query("name"), c.Query("section"), c.Query("status"),
		helper.ParseInt64(c.Query("departmentId"), -1),
		helper.ParseInt64(c.Query("cursor"), -1),
		helper.ParseInt(c.Query("limit"), 20),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": operators})
}

// GetOperatorsByDepartment lists a department's operators.
func GetOperatorsByDepartment(c *gin.Context) {
	deptID, err := strconv.ParseInt(c.Param("departmentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department id"})
		return
	}
	op, err := operatorSvc.GetOperatorsByDepartment(deptID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": op})
}

// FindOperatorsByShiftDept lists operators for a department + shift.
func FindOperatorsByShiftDept(c *gin.Context) {
	op, err := operatorSvc.FindOperatorsByShiftDept(
		helper.ParseInt64(c.Query("departmentId"), -1),
		helper.ParseInt64(c.Query("shiftId"), -1),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": op})
}

// CreateOperators bulk-creates operators from a JSON array (or a single object).
func CreateOperators(c *gin.Context) {
	var payload json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	result, err := operatorSvc.CreateOperators(payload, helper.ResolveActor(c, helper.ActorFromJSON(payload)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Operators created successfully", "data": result})
}

// UpdateOperator patches an operator.
func UpdateOperator(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operator id"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	operatorID, err := operatorSvc.UpdateOperator(id, updates, helper.ResolveActor(c, helper.ActorFromMap(updates)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if operatorID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Operator not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Operator updated successfully", "operatorId": operatorID})
}

// GetShifts lists all shifts.
func GetShifts(c *gin.Context) {
	shifts, err := shiftSvc.GetShifts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shifts})
}

// GetShift returns one shift by id.
func GetShift(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shift id"})
		return
	}
	s, err := shiftSvc.GetShiftById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(s) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shift not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": s})
}

// GetShiftsByDepartment lists a department's shifts.
func GetShiftsByDepartment(c *gin.Context) {
	deptID, err := strconv.Atoi(c.Param("departmentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department id"})
		return
	}
	s, err := shiftSvc.GetShiftsByDepartment(deptID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": s})
}

// GetShiftOutputHistory returns one shift's output history.
func GetShiftOutputHistory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shift id"})
		return
	}
	rows, err := shiftSvc.GetShiftOutputHistory(id, c.Query("from"), c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

// FindShiftOutputHistory lists output history by department/shift/date.
func FindShiftOutputHistory(c *gin.Context) {
	rows, err := shiftSvc.FindShiftOutputHistory(
		helper.ParseInt64(c.Query("departmentId"), -1),
		helper.ParseInt64(c.Query("shiftId"), -1),
		c.Query("from"), c.Query("to"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

// CreateShifts bulk-creates shifts from a JSON array (or a single object).
func CreateShifts(c *gin.Context) {
	var payload json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	result, err := shiftSvc.CreateShifts(payload, helper.ResolveActor(c, helper.ActorFromJSON(payload)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Shifts created successfully", "data": result})
}

// UpdateShift patches a shift.
func UpdateShift(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shift id"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	shiftID, err := shiftSvc.UpdateShift(id, updates, helper.ResolveActor(c, helper.ActorFromMap(updates)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if shiftID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shift not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shift updated successfully", "shiftId": shiftID})
}

// GetDevice returns one device by id.
func GetDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device id"})
		return
	}
	d, err := deviceSvc.GetDeviceById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(d) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": d})
}

// FindDevices lists devices by filters with cursor pagination.
func FindDevices(c *gin.Context) {
	devices, err := deviceSvc.FindDevices(
		c.Query("androidId"),
		helper.ParseInt64(c.Query("processId"), -1),
		c.Query("status"),
		helper.ParseInt64(c.Query("machineId"), -1),
		helper.ParseInt64(c.Query("cursor"), -1),
		helper.ParseInt(c.Query("limit"), 20),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": devices})
}

// IdentifyDevice heartbeats a device by android id and returns its context.
func IdentifyDevice(c *gin.Context) {
	androidID := c.Param("androidId")
	if androidID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "androidId is required"})
		return
	}
	d, err := deviceSvc.IdentifyDevice(androidID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(d) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": d})
}

// CreateDevices bulk-creates devices from a JSON array.
func CreateDevices(c *gin.Context) {
	var payload json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	result, err := deviceSvc.CreateDevices(payload, helper.ResolveActor(c, helper.ActorFromJSON(payload)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Devices created successfully", "data": result})
}

// UpdateDevice patches a device.
func UpdateDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device id"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	deviceID, err := deviceSvc.UpdateDevice(id, updates, helper.ResolveActor(c, helper.ActorFromMap(updates)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if deviceID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Device updated successfully", "deviceId": deviceID})
}

// GetMachineType returns one machine type by id, including nested machines.
func GetMachineType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine type id"})
		return
	}
	mt, err := machineTypeSvc.GetMachineTypeById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(mt) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Machine type not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": mt})
}

// CreateMachineTypes bulk-creates machine types (and nested machines) from a JSON array (or a single object).
func CreateMachineTypes(c *gin.Context) {
	var payload json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	result, err := machineTypeSvc.CreateMachineTypes(payload, helper.ResolveActor(c, helper.ActorFromJSON(payload)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Machine types created successfully", "data": result})
}

// GetMachine returns one machine by id, including unit, hold reasons, performance, and current operation.
func GetMachine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine id"})
		return
	}
	m, err := machineSvc.GetMachineById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(m) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Machine not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

// FindMachines lists machines by filters with cursor pagination.
func FindMachines(c *gin.Context) {
	machines, err := machineSvc.FindMachines(
		c.Query("code"),
		c.Query("name"),
		helper.ParseInt64(c.Query("machineTypeId"), -1),
		helper.ParseInt64(c.Query("departmentId"), -1),
		helper.ParseInt64(c.Query("processId"), -1),
		helper.ParseInt64(c.Query("cursor"), -1),
		helper.ParseInt(c.Query("limit"), 20),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": machines})
}

// CreateMachines bulk-creates machines from a JSON array (or a single object).
func CreateMachines(c *gin.Context) {
	var payload json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	result, err := machineSvc.CreateMachines(payload, helper.ResolveActor(c, helper.ActorFromJSON(payload)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Machines created successfully", "data": result})
}

// UpdateMachine patches a machine.
func UpdateMachine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine id"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	machineID, err := machineSvc.UpdateMachine(id, updates, helper.ResolveActor(c, helper.ActorFromMap(updates)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if machineID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Machine not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Machine updated successfully", "machineId": machineID})
}
