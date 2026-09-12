package auth

import "time"

type User struct {
	ID           int64
	Name         string
	Email        string
	HashPassword string
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}
