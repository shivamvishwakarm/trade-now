package http

import "github.com/gin-gonic/gin"

func registerWebSocketRoutes(router *gin.Engine, deps RouterDeps) {

	router.GET("/ws", deps.WebSocketHandler.Handle)

}
