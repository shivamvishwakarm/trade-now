package user

import (
	"context"
	"errors"

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
func (s *Service) GetProfile(ctx context.Context) (Profile, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		s.logger.Error("failed to get auth claims from context")
		return Profile{}, ErrInvalidToken
	}

	id := claims.ID

	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Error(
			"failed to get user by ID",
			zap.Error(err),
		)

		return Profile{}, err
	}

	return Profile{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Name,
	}, nil
}
