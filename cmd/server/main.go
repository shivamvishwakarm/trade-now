package main

import (
	"fmt"
	"log"

	"github.com/shivamvishwakarm/trade-now/internal/auth"
	"github.com/shivamvishwakarm/trade-now/internal/config"
	"github.com/shivamvishwakarm/trade-now/internal/database"
	"github.com/shivamvishwakarm/trade-now/internal/http"
	"github.com/shivamvishwakarm/trade-now/internal/websocket"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	dbConfig, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Database connection
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.DB.Host,
		dbConfig.DB.Port,
		dbConfig.DB.User,
		dbConfig.DB.Password,
		dbConfig.DB.Name,
	)

	db, err := database.NewPostgres(database.Config{
		DSN: dsn,
	})
	if err != nil {
		logger.Fatal("database connection failed", zap.Error(err))
	}
	defer db.Close()

	// Dependencies
	passwordHasher := auth.NewBcryptPasswordHasher(
		bcrypt.DefaultCost,
	)

	userRepository := database.NewUserRepository(db)

	authService := auth.NewService(auth.ServiceDeps{
		Logger:         logger,
		UserRepo:       userRepository,
		PasswordHasher: passwordHasher,
	})

	// Handlers
	websocketHandler := websocket.NewHandler(
		websocket.HandlerDeps{
			Logger: logger,
		},
	)

	authHandler := auth.NewHandler(
		auth.HandlerDeps{
			Logger:  logger,
			Service: authService,
		},
	)

	// Router
	router := http.NewRouter(http.RouterDeps{
		WebSocketHandler: websocketHandler,
		AuthHandler:      authHandler,
	})

	logger.Info(
		"starting server",
		zap.String("addr", ":8080"),
	)

	if err := router.Run(":8080"); err != nil {
		logger.Fatal(
			"server stopped",
			zap.Error(err),
		)
	}
}
