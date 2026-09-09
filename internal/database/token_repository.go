package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/shivamvishwakarm/trade-now/internal/auth"
)

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) StoreRefreshToken(
	ctx context.Context,
	userID int64,
	tokenHash string,
	expiresAt time.Time,
) error {
	const query = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (r *TokenRepository) GetRefreshToken(
	ctx context.Context,
	tokenHash string,
) (auth.RefreshToken, error) {
	const query = `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	var t auth.RefreshToken
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&t.ID,
		&t.UserID,
		&t.TokenHash,
		&t.ExpiresAt,
		&t.CreatedAt,
	)
	if err != nil {
		return auth.RefreshToken{}, err
	}
	return t, nil
}

func (r *TokenRepository) DeleteRefreshToken(
	ctx context.Context,
	tokenHash string,
) error {
	const query = `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (r *TokenRepository) DeleteAllUserTokens(
	ctx context.Context,
	userID int64,
) error {
	const query = `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
