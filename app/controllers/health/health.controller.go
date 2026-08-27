package health

import (
	"github.com/gin-gonic/gin"

	healthService "github.com/thimira/production-tracer/app/services/health"
)

// Init registers the health-check routes on the router.
func Init(router *gin.Engine) {
	router.GET("/health", healthService.HealthCheck)
}
