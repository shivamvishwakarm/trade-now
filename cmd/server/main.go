package main

import (
	"log"

	"github.com/shivamvishwakarm/trade-now/internal/auth"
	"github.com/shivamvishwakarm/trade-now/internal/http"
	"github.com/shivamvishwakarm/trade-now/internal/websocket"
	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("fail to initialize logger: %v", err)
	}

	defer logger.Sync()

	websocketHandler := websocket.NewHandler(websocket.HandlerDeps{
		Logger: logger,
	})

	authHandler := auth.NewHandler(auth.HandlerDeps{
		Logger: logger,
	})

	router := http.NewRouter(http.RouterDeps{
		WebSocketHandler: websocketHandler,
		AuthHandler:      authHandler,
	})

	logger.Info("starting server", zap.String("addr", ":8080"))

	if err := router.Run(":8080"); err != nil {
		logger.Fatal("server stopped", zap.Error(err))
	}
}
