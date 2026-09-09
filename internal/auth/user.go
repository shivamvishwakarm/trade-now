package auth

type User struct {
	ID           int64
	Name         string
	Email        string
	HashPassword string
}
