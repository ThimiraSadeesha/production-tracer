package api

import (
	"github.com/thimira/production-tracer/app/controllers"
	corsUtil "github.com/thimira/production-tracer/internal/cors"
	"github.com/thimira/production-tracer/internal/interceptors"
	"github.com/thimira/production-tracer/internal/middleware"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var router *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	router.SetTrustedProxies(nil)
	router.Use(middleware.ErrorHandler())
	router.Use(corsUtil.CORS())
	router.Use(requestid.New())
	router.Use(interceptors.Interceptor())
	router.Use(interceptors.Log())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	controllers.Init(router)
}

func Run(addr string) {
	if err := router.Run(addr); err != nil {
		panic(err)
	}
}

func Router() *gin.Engine {
	return router
}
