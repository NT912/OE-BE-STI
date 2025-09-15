package routes

import (
	v1 "gateway/api/v1"
	"gateway/api/v1/auth"
	lauchpad "gateway/api/v1/launchpad"
	"gateway/api/v1/users"
	wallets "gateway/api/v1/wallet"
	"gateway/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.ExceptionMiddleware())
	// r.Use(middlewares.ResponseFormatter())
	r.GET("/api/v1/health", v1.HealthCheck)
	wallets.InitModule()
	users.InitModule(r)
	auth.InitModule(r)
	lauchpad.InitModule(r)

	return r
}
