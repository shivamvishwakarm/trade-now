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

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Database connection
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)

	db, err := database.NewPostgres(database.Config{DSN: dsn})
	if err != nil {
		logger.Fatal("database connection failed", zap.Error(err))
	}
	defer db.Close()

	// Repositories
	userRepository := database.NewUserRepository(db)
	tokenRepository := database.NewTokenRepository(db)

	// Dependencies
	passwordHasher := auth.NewBcryptPasswordHasher(bcrypt.DefaultCost)

	authService := auth.NewService(auth.ServiceDeps{
		Logger:         logger,
		UserRepo:       userRepository,
		TokenRepo:      tokenRepository,
		PasswordHasher: passwordHasher,
		AccessSecret:   cfg.JWT.AccessSecret,
		RefreshSecret:  cfg.JWT.RefreshSecret,
		AccessExpiry:   cfg.JWT.AccessTokenExpiry,
		RefreshExpiry:  cfg.JWT.RefreshTokenExpiry,
	})

	// Handlers
	websocketHandler := websocket.NewHandler(websocket.HandlerDeps{Logger: logger})
	authHandler := auth.NewHandler(auth.HandlerDeps{
		Logger:  logger,
		Service: authService,
	})

	// Router
	router := http.NewRouter(http.RouterDeps{
		WebSocketHandler: websocketHandler,
		AuthHandler:      authHandler,
	})

	logger.Info("starting server", zap.String("addr", ":8080"))

	if err := router.Run(":8080"); err != nil {
		logger.Fatal("server stopped", zap.Error(err))
	}
}
