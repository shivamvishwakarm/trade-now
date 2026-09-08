package auth

type PasswordHasher interface {
	Hash(password string) (string, error)
}
