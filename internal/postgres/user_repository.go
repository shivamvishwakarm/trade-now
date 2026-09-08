package postgres

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

func (r *UserRepository) ExistByEmail(ctx context.Context, email string) (bool, error) {

	const query = `SELECT EXISTS (
		select 1 from users WHERE email = $1
	)`

	var exists bool

	if err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(&exists); err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user auth.User) error {

	const query = `INSERT INTO users (
		id,email,password_hash
	) VALUES ($1,$2,$3)
		`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
	)
	return err

}
