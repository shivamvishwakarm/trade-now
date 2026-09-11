package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shivamvishwakarm/trade-now/internal/auth"
	"github.com/shivamvishwakarm/trade-now/internal/user"
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
	var u auth.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.HashPassword,
	)

	if err != nil {
		return auth.User{}, err
	}

	return u, nil
}

// UserProfileRepository implements user.UserRepository, returning the public
// user.User type (no sensitive fields like password hash).
type UserProfileRepository struct {
	db *sql.DB
}

func NewUserProfileRepository(db *sql.DB) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

func (r *UserProfileRepository) GetByEmail(
	ctx context.Context,
	email string,
) (user.Profile, error) {
	const query = `SELECT id, name, email FROM users WHERE email = $1`

	var u user.Profile
	var id int64

	err := r.db.QueryRowContext(ctx, query, email).Scan(&id, &u.Name, &u.Email)
	if err != nil {
		return user.Profile{}, err
	}

	u.Id = fmt.Sprintf("%d", id)
	return u, nil
}
