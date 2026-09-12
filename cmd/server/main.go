// @title           Trade Now API
// @version         1.0
// @description     REST API for the Trade Now application. Handles authentication and real-time WebSocket connectivity.

// @contact.name    Trade Now Support
// @contact.email   support@trade-now.dev

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey  AccessTokenCookie
// @in                          cookie
// @name                        access_token
// @description                 JWT access token stored as an HttpOnly cookie.

// @securityDefinitions.apikey  RefreshTokenCookie
// @in                          cookie
// @name                        refresh_token
// @description                 JWT refresh token stored as an HttpOnly cookie (only sent to /auth/refresh).

package main

import (
	"fmt"
	"log"

	_ "github.com/shivamvishwakarm/trade-now/docs"
	"github.com/shivamvishwakarm/trade-now/internal/auth"
	"github.com/shivamvishwakarm/trade-now/internal/config"
	"github.com/shivamvishwakarm/trade-now/internal/database"
	"github.com/shivamvishwakarm/trade-now/internal/http"
	"github.com/shivamvishwakarm/trade-now/internal/instrument"
	"github.com/shivamvishwakarm/trade-now/internal/user"
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
	userProfileRepository := database.NewUserProfileRepository(db)
	tokenRepository := database.NewTokenRepository(db)
	instrumentRepository := database.NewInstrumentRepository(db)

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

	userService := user.NewService(user.ServiceDeps{
		Logger:   logger,
		UserRepo: userProfileRepository,
	})

	instrumentService := instrument.NewService(instrument.ServiceDeps{
		Logger:         logger,
		InstrumentRepo: instrumentRepository,
	})

	// Handlers
	websocketHandler := websocket.NewHandler(websocket.HandlerDeps{Logger: logger})
	authHandler := auth.NewHandler(auth.HandlerDeps{
		Logger:  logger,
		Service: authService,
	})

	userHandler := user.NewHandler(user.HandlerDeps{
		Logger:  logger,
		Service: userService,
	})

	instrumentHandler := instrument.NewHandler(instrument.HandlerDeps{
		Logger:  logger,
		Service: instrumentService,
	})

	// Router
	router := http.NewRouter(http.RouterDeps{
		WebSocketHandler:  websocketHandler,
		AuthHandler:       authHandler,
		UserHandler:       userHandler,
		InstrumentHandler: instrumentHandler,
		AccessSecret:      cfg.JWT.AccessSecret,
	})

	logger.Info("starting server", zap.String("addr", ":8080"))

	if err := router.Run(":8080"); err != nil {
		logger.Fatal("server stopped", zap.Error(err))
	}
}
