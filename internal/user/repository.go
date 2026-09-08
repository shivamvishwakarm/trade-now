package user

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
}

type UserRepository interface {
	ExistByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user User) error
}
