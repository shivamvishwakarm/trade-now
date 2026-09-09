package auth

import (
	"context"
	"time"
)

type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user User) error
	GetByEmail(ctx context.Context, email string) (User, error)
}

type TokenRepository interface {
	// StoreRefreshToken persists a hashed refresh token for the given user.
	StoreRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	// GetRefreshToken looks up a stored token by its hash.
	GetRefreshToken(ctx context.Context, tokenHash string) (RefreshToken, error)
	// DeleteRefreshToken removes a single token (logout).
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	// DeleteAllUserTokens removes every token for a user (force-logout all sessions).
	DeleteAllUserTokens(ctx context.Context, userID int64) error
}
