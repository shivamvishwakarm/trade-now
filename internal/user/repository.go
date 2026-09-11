package user

import "context"

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (Profile, error)
	GetById(ctx context.Context, id string) (Profile, error)
}
