package controllers

import (
	"github.com/gin-gonic/gin"

	healthController "github.com/thimira/production-tracer/app/controllers/health"
	setupController "github.com/thimira/production-tracer/app/controllers/setup"
	workOrderController "github.com/thimira/production-tracer/app/controllers/work-order"
)

// Init wires every feature controller onto the router. Register new controllers
// here as they are added.
func Init(router *gin.Engine) {
	healthController.Init(router)
	setupController.Init(router)
	workOrderController.Init(router)
}
