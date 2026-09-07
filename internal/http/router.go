package http

import (
	"net/http"

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

	router.GET("/ping", func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Pong",
		})
	})

	router.GET("/connect", deps.WsHandler.Handle)

	return router

}
