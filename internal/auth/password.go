package auth

import (
	"strings"
	"unicode"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, plainPassword string) error
}

func PasswordStrengthCheck(password string) (string, string) {
	password = strings.TrimSpace(password)

	score := 0

	// At least 8 characters
	if len(password) >= 8 {
		score++
	}

	// Uppercase letter
	if strings.IndexFunc(password, unicode.IsUpper) != -1 {
		score++
	}

	// Lowercase letter
	if strings.IndexFunc(password, unicode.IsLower) != -1 {
		score++
	}

	// Number
	if strings.IndexFunc(password, unicode.IsDigit) != -1 {
		score++
	}

	// Special character
	if strings.IndexFunc(password, func(r rune) bool {
		return strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:',.<>?/`~", r)
	}) != -1 {
		score++
	}

	var strength string

	switch {
	case score <= 2:
		strength = "Weak"
	case score <= 4:
		strength = "Medium"
	default:
		strength = "Strong"
	}

	return password, strength
}
