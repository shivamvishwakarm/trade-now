package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamvishwakarm/trade-now/internal/websocket"
)

type RouterDeps struct {
	WebSocketHandler *websocket.Handler
}

func NewRouter(deps RouterDeps) *gin.Engine {

	router := gin.New()

	router.Use(gin.Recovery())

	registerHealthRoute(router)

	registerWebsocketRoute(router, deps)

	return router

}
