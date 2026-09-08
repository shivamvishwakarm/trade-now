package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shivamvishwakarm/trade-now/internal/websocket"
	"go.uber.org/zap"
)

type RouterDeps struct {
	WsHandler *websocket.Handler
	Logger    *zap.Logger
}

func NewRouter(deps RouterDeps) *gin.Engine {

	router := gin.New()

	router.Use(gin.Recovery())

	routesV1 := router.Group("/api/v1")

	routesV1.GET("/health", func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, gin.H{
			"message":   "healthy",
			"timestamp": time.Now().UTC(),
		})
	})

	routesV1.GET("/connect", deps.WsHandler.Handle)

	return router

}
