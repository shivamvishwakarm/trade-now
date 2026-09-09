package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
)

var (
	ErrInvalidEmail    = errors.New("email is required")
	ErrInvalidPassword = errors.New("password is required")
	ErrEmailExists     = errors.New("user with email already exists")
	ErrUserNotFound    = errors.New("invalid email or password")
)

type Service struct {
	logger         *zap.Logger
	passwordHasher PasswordHasher
	userRepo       UserRepository
	tokenRepo      TokenRepository
	accessSecret   string
	refreshSecret  string
	accessExpiry   time.Duration
	refreshExpiry  time.Duration
}

type ServiceDeps struct {
	Logger         *zap.Logger
	PasswordHasher PasswordHasher
	UserRepo       UserRepository
	TokenRepo      TokenRepository
	AccessSecret   string
	RefreshSecret  string
	AccessExpiry   time.Duration
	RefreshExpiry  time.Duration
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		logger:         deps.Logger,
		passwordHasher: deps.PasswordHasher,
		userRepo:       deps.UserRepo,
		tokenRepo:      deps.TokenRepo,
		accessSecret:   deps.AccessSecret,
		refreshSecret:  deps.RefreshSecret,
		accessExpiry:   deps.AccessExpiry,
		refreshExpiry:  deps.RefreshExpiry,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) error {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if email == "" {
		return ErrInvalidEmail
	}
	if req.Password == "" {
		return ErrInvalidPassword
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to check whether user exists", zap.Error(err))
		return err
	}
	if exists {
		return ErrEmailExists
	}

	passwordHash, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))
		return err
	}

	user := User{
		Name:         req.Name,
		Email:        email,
		HashPassword: passwordHash,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (User, TokenPair, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if email == "" {
		return User{}, TokenPair{}, ErrInvalidEmail
	}
	if req.Password == "" {
		return User{}, TokenPair{}, ErrInvalidPassword
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, TokenPair{}, ErrUserNotFound
		}
		s.logger.Error("failed to get user by email", zap.Error(err))
		return User{}, TokenPair{}, err
	}

	if err := s.passwordHasher.Compare(user.HashPassword, req.Password); err != nil {
		return User{}, TokenPair{}, ErrUserNotFound
	}

	pair, err := GenerateTokenPair(
		user.ID, user.Email,
		s.accessSecret, s.refreshSecret,
		s.accessExpiry, s.refreshExpiry,
	)
	if err != nil {
		s.logger.Error("failed to generate token pair", zap.Error(err))
		return User{}, TokenPair{}, err
	}

	tokenHash := HashToken(pair.RefreshToken)
	expiresAt := time.Now().Add(s.refreshExpiry)

	if err := s.tokenRepo.StoreRefreshToken(ctx, user.ID, tokenHash, expiresAt); err != nil {
		s.logger.Error("failed to store refresh token", zap.Error(err))
		return User{}, TokenPair{}, err
	}

	return user, pair, nil
}

func (s *Service) Refresh(ctx context.Context, rawRefreshToken string) (TokenPair, error) {
	// Validate the incoming refresh token signature + expiry
	claims, err := ParseToken(rawRefreshToken, s.refreshSecret)
	if err != nil {
		return TokenPair{}, ErrInvalidToken
	}

	// Verify it exists in the DB (not revoked)
	tokenHash := HashToken(rawRefreshToken)
	stored, err := s.tokenRepo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TokenPair{}, ErrInvalidToken
		}
		s.logger.Error("failed to look up refresh token", zap.Error(err))
		return TokenPair{}, err
	}

	if time.Now().After(stored.ExpiresAt) {
		_ = s.tokenRepo.DeleteRefreshToken(ctx, tokenHash)
		return TokenPair{}, ErrInvalidToken
	}

	// Rotate: delete old, issue new pair
	if err := s.tokenRepo.DeleteRefreshToken(ctx, tokenHash); err != nil {
		s.logger.Error("failed to delete old refresh token", zap.Error(err))
		return TokenPair{}, err
	}

	pair, err := GenerateTokenPair(
		claims.UserID, claims.Email,
		s.accessSecret, s.refreshSecret,
		s.accessExpiry, s.refreshExpiry,
	)
	if err != nil {
		s.logger.Error("failed to generate token pair", zap.Error(err))
		return TokenPair{}, err
	}

	newHash := HashToken(pair.RefreshToken)
	if err := s.tokenRepo.StoreRefreshToken(ctx, claims.UserID, newHash, time.Now().Add(s.refreshExpiry)); err != nil {
		s.logger.Error("failed to store new refresh token", zap.Error(err))
		return TokenPair{}, err
	}

	return pair, nil
}

func (s *Service) Logout(ctx context.Context, rawRefreshToken string) error {
	tokenHash := HashToken(rawRefreshToken)
	if err := s.tokenRepo.DeleteRefreshToken(ctx, tokenHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil // already gone, treat as success
		}
		s.logger.Error("failed to delete refresh token on logout", zap.Error(err))
		return err
	}
	return nil
}
