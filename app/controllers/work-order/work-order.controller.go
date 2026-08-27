package workorder

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	workOrderService "github.com/thimira/production-tracer/app/services/work-order"
)

func Init(router *gin.Engine) {
	group := router.Group("/work-orders")
	{
		group.GET("/find", FindWorkOrders)
		group.GET("/:id", GetWorkOrder)
		group.POST("", workOrderService.Save)
		group.PATCH("/:id", UpdateWorkOrder)
	}
}

// GetWorkOrder returns one work order with its items.
func GetWorkOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work order id"})
		return
	}
	wo, err := workOrderService.GetWorkOrderById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(wo) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": wo})
}

// FindWorkOrders lists work orders by filters with cursor pagination.
func FindWorkOrders(c *gin.Context) {
	cursor, _ := strconv.ParseInt(c.Query("cursor"), 10, 64)
	limit, _ := strconv.Atoi(c.Query("limit"))
	list, err := workOrderService.FindWorkOrders(c.Query("reference"), c.Query("customer"), c.Query("status"), cursor, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// UpdateWorkOrder patches a work order header.
func UpdateWorkOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work order id"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	woID, err := workOrderService.UpdateWorkOrder(id, updates, "system")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if woID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Work order updated successfully", "workOrderId": woID})
}
