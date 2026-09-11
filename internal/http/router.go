package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamvishwakarm/trade-now/internal/auth"
	"github.com/shivamvishwakarm/trade-now/internal/user"
	"github.com/shivamvishwakarm/trade-now/internal/websocket"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterDeps struct {
	WebSocketHandler *websocket.Handler
	AuthHandler      *auth.Handler
	UserHandler      *user.Handler
	AccessSecret     string
	// UserHandler *user.Handler  — uncomment once the user handler is implemented
}

func NewRouter(deps RouterDeps) *gin.Engine {

	router := gin.New()

	router.Use(gin.Recovery())

	routerV1 := router.Group("/api/v1")

	// Swagger UI — available at /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	registerHealthRoutes(router)

	registerWebSocketRoutes(router, deps)

	registerAuthRoutes(routerV1, deps)

	// Protected routes — require a valid access_token cookie
	protected := routerV1.Group("", RequireAuth(deps.AccessSecret))
	userRoutes(protected, deps)

	return router
}
