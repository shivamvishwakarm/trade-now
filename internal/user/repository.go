package user

import (
	"context"
)

type User struct {
	Email        string
	PasswordHash string
}

type UserRepository interface {
	ExistByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, user User) error
}
