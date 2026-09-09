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
			email,
			name,
			password_hash
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		u.Email,
		u.Name,
		u.HashPassword,
	)

	return err
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (auth.User, error) {

	const query = `
		SELECT id, name, email, password_hash FROM users WHERE email = $1
	`
	var user auth.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashPassword,
	)

	if err != nil {
		return auth.User{}, err
	}

	return user, nil
}
