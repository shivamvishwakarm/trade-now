package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	logger         *zap.Logger
	passwordHasher PasswordHasher
	userRepo       UserRepository
}

type ServiceDeps struct {
	Logger         *zap.Logger
	PasswordHasher PasswordHasher
	UserRepo       UserRepository
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		logger:         deps.Logger,
		passwordHasher: deps.PasswordHasher,
		userRepo:       deps.UserRepo,
	}
}

var (
	ErrInvalidEmail    = errors.New("email is required")
	ErrInvalidPassword = errors.New("password is required")
	ErrEmailExists     = errors.New("user with email already exists")
)

func (s *Service) Register(
	ctx context.Context,
	req RegisterRequest,
) error {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if email == "" {
		return ErrInvalidEmail
	}

	if req.Password == "" {
		return ErrInvalidPassword
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		s.logger.Error(
			"failed to check whether user exists",
			zap.Error(err),
		)
		return err
	}

	if exists {
		return ErrEmailExists
	}

	passwordHash, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		s.logger.Error(
			"failed to hash password",
			zap.Error(err),
		)
		return err
	}

	user := User{
		Id:           uuid.New(),
		Email:        email,
		HashPassword: passwordHash,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error(
			"failed to create user",
			zap.Error(err),
		)
		return err
	}

	return nil
}
