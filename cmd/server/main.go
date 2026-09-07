package main

import (
	"log"

	"github.com/shivamvishwakarm/trade-now/internal/http"
	"github.com/shivamvishwakarm/trade-now/internal/websocket"
	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("fail to initialize logger: %w", err)
	}

	defer logger.Sync()

	wsHandler := websocket.NewHandler(websocket.HandlerDeps{
		Logger: logger,
	})

	router := http.NewRouter(http.RouterDeps{
		WsHandler: wsHandler,
		Logger:    logger,
	})

	logger.Info("starting server", zap.String("addr", ":8080"))

	if err := router.Run(":8080"); err != nil {
		logger.Fatal("server stopped", zap.Error(err))
	}
}
