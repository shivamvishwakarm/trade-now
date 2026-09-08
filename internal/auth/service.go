package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
}

type ServiceDeps struct {
	Logger         *zap.Logger
	PasswordHasher PasswordHasher
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		logger: deps.Logger,
	}
}

func (s *Service) Register(ctx *gin.Context, req RegisterRequest) error {

	// 1. Normalize email
	// 2. Validate business rules
	// 3. Hash password
	// 4. Build User
	// 5. Persist User

	return nil
}
