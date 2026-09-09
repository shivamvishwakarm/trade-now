package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamvishwakarm/trade-now/internal/auth"
	"github.com/shivamvishwakarm/trade-now/internal/websocket"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterDeps struct {
	WebSocketHandler *websocket.Handler
	AuthHandler      *auth.Handler
}

func NewRouter(deps RouterDeps) *gin.Engine {

	router := gin.New()

	router.Use(gin.Recovery())

	// Swagger UI — available at /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	registerHealthRoutes(router)

	registerWebSocketRoutes(router, deps)

	registerAuthRoutes(router, deps)

	return router
}
