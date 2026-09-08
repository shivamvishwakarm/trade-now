package user

import (
	"context"

	"github.com/shivamvishwakarm/trade-now/internal/auth"
)

type UserRepository interface {
	ExistByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user auth.User) error
}
