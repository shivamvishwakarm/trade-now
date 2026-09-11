package http

import "github.com/gin-gonic/gin"

func userRoutes(router *gin.RouterGroup, deps RouterDeps) {
	me := router.Group("/me")

	me.GET("/", deps.UserHandler.GetMe)

}
