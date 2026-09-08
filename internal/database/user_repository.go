package database

import (
	"context"
	"database/sql"

	"github.com/shivamvishwakarm/trade-now/internal/auth"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) ExistsByEmail(
	ctx context.Context,
	email string,
) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
		)
	`

	var exists bool

	if err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) Create(
	ctx context.Context,
	u auth.User,
) error {
	const query = `
		INSERT INTO users (
			id,
			email,
			password_hash
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		u.Id,
		u.Email,
		u.HashPassword,
	)

	return err
}
