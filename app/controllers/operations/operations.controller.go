// Package operations exposes the production execution routes (operation processes).
package operations

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	operationSvc "github.com/thimira/production-tracer/app/services/operation-process"
	"github.com/thimira/production-tracer/internal/helper"
)

func Init(router *gin.Engine) {
	group := router.Group("/operations")
	{
		group.GET("/find", FindOperations)
		group.GET("/:id", GetOperation)
		group.POST("", CreateOperations)
		group.PATCH("/:id", UpdateOperation)
	}
}

// CreateOperations bulk-creates operation processes from a JSON array.
func CreateOperations(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	raw := json.RawMessage(body)
	result, err := operationSvc.Save(raw, helper.ResolveActor(c, helper.ActorFromJSON(raw)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Operation processes created successfully", "data": result})
}

// GetOperation returns one operation process with its attributes.
func GetOperation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation id"})
		return
	}
	op, err := operationSvc.GetOperationProcessById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(op) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Operation not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": op})
}

// FindOperations lists operation processes by filters with cursor pagination.
func FindOperations(c *gin.Context) {
	list, err := operationSvc.FindOperationProcesses(
		helper.ParseInt64(c.Query("machineId"), -1),
		helper.ParseInt64(c.Query("processId"), -1),
		c.Query("plannedDate"),
		c.Query("workOrderId"),
		c.Query("status"),
		helper.ParseInt64(c.Query("cursor"), -1),
		helper.ParseInt(c.Query("limit"), 20),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// UpdateOperation patches an operation process.
func UpdateOperation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation id"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	opID, err := operationSvc.UpdateOperationProcess(id, updates, helper.ResolveActor(c, helper.ActorFromMap(updates)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if opID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Operation not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Operation updated successfully", "operationProcessId": opID})
}
