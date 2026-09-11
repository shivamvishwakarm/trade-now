package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/shivamvishwakarm/trade-now/internal/auth"

	"go.uber.org/zap"
)

type Service struct {
	logger   *zap.Logger
	userRepo UserRepository
}

type ServiceDeps struct {
	Logger   *zap.Logger
	UserRepo UserRepository
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		logger:   deps.Logger,
		userRepo: deps.UserRepo,
	}
}

var (
	ErrInvalidToken = errors.New("invalid auth token")
	ErrUserNotExist = errors.New("User not exists")
)

// GetProfile returns the profile for the authenticated user.
// TODO: implement
func (s *Service) GetProfile(ctx *gin.Context) (Profile, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		s.logger.Error("Can't get claims from context")
		return Profile{}, ErrInvalidToken
	}

	email := claims.Email

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to get user by email")
		return Profile{}, ErrUserNotExist
	}

	return Profile{
		Name:  user.Name,
		Email: user.Name,
		Id:    user.Id,
	}, nil
}
