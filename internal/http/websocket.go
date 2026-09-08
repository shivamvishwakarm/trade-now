package http

import "github.com/gin-gonic/gin"

func registerWebsocketRoute(router *gin.Engine, deps RouterDeps) {

	router.GET("/ws", deps.WebSocketHandler.Handle)

}
