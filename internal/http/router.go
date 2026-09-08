package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamvishwakarm/trade-now/internal/auth"
	"github.com/shivamvishwakarm/trade-now/internal/websocket"
)

type RouterDeps struct {
	WebSocketHandler *websocket.Handler
	AuthHandler      *auth.Handler
}

func NewRouter(deps RouterDeps) *gin.Engine {

	router := gin.New()

	router.Use(gin.Recovery())

	registerHealthRoutes(router)

	registerWebSocketRoutes(router, deps)

	registerAuthRoutes(router, deps)

	return router
}
