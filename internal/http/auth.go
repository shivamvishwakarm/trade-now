package http

import (
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(router *gin.RouterGroup, deps RouterDeps) {

	authRouter := router.Group("/auth")
	authRouter.POST("/register", deps.AuthHandler.Register)
	authRouter.POST("/login", deps.AuthHandler.Login)
	authRouter.POST("refresh", deps.AuthHandler.Refresh)
	authRouter.POST("/logout", deps.AuthHandler.Logout)

}
