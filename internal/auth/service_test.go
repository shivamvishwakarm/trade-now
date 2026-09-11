package auth

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
)

// ---------------------------------------------------------------------------
// Mock implementations
// ---------------------------------------------------------------------------

type mockUserRepo struct {
	existsByEmailFn func(ctx context.Context, email string) (bool, error)
	createFn        func(ctx context.Context, user User) error
	getByEmailFn    func(ctx context.Context, email string) (User, error)
}

func (m *mockUserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return m.existsByEmailFn(ctx, email)
}
func (m *mockUserRepo) Create(ctx context.Context, user User) error {
	return m.createFn(ctx, user)
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (User, error) {
	return m.getByEmailFn(ctx, email)
}

type mockTokenRepo struct {
	storeRefreshTokenFn   func(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	getRefreshTokenFn     func(ctx context.Context, tokenHash string) (RefreshToken, error)
	deleteRefreshTokenFn  func(ctx context.Context, tokenHash string) error
	deleteAllUserTokensFn func(ctx context.Context, userID int64) error
}

func (m *mockTokenRepo) StoreRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	return m.storeRefreshTokenFn(ctx, userID, tokenHash, expiresAt)
}
func (m *mockTokenRepo) GetRefreshToken(ctx context.Context, tokenHash string) (RefreshToken, error) {
	return m.getRefreshTokenFn(ctx, tokenHash)
}
func (m *mockTokenRepo) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	return m.deleteRefreshTokenFn(ctx, tokenHash)
}
func (m *mockTokenRepo) DeleteAllUserTokens(ctx context.Context, userID int64) error {
	return m.deleteAllUserTokensFn(ctx, userID)
}

type mockPasswordHasher struct {
	hashFn    func(password string) (string, error)
	compareFn func(hashed, plain string) error
}

func (m *mockPasswordHasher) Hash(password string) (string, error) {
	return m.hashFn(password)
}
func (m *mockPasswordHasher) Compare(hashed, plain string) error {
	return m.compareFn(hashed, plain)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestService(userRepo UserRepository, tokenRepo TokenRepository, hasher PasswordHasher) *Service {
	return NewService(ServiceDeps{
		Logger:         zap.NewNop(),
		PasswordHasher: hasher,
		UserRepo:       userRepo,
		TokenRepo:      tokenRepo,
		AccessSecret:   "access-secret",
		RefreshSecret:  "refresh-secret",
		AccessExpiry:   15 * time.Minute,
		RefreshExpiry:  7 * 24 * time.Hour,
	})
}

// ---------------------------------------------------------------------------
// Register tests
// ---------------------------------------------------------------------------

func TestService_Register(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userRepo := &mockUserRepo{
			existsByEmailFn: func(_ context.Context, _ string) (bool, error) { return false, nil },
			createFn:        func(_ context.Context, _ User) error { return nil },
		}
		hasher := &mockPasswordHasher{
			hashFn: func(_ string) (string, error) { return "hashed", nil },
		}
		svc := newTestService(userRepo, &mockTokenRepo{}, hasher)

		err := svc.Register(ctx, RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "secret"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("empty email returns ErrInvalidEmail", func(t *testing.T) {
		svc := newTestService(&mockUserRepo{}, &mockTokenRepo{}, &mockPasswordHasher{})

		err := svc.Register(ctx, RegisterRequest{Email: "   ", Password: "secret"})
		if !errors.Is(err, ErrInvalidEmail) {
			t.Fatalf("expected ErrInvalidEmail, got %v", err)
		}
	})

	t.Run("empty password returns ErrInvalidPassword", func(t *testing.T) {
		svc := newTestService(&mockUserRepo{}, &mockTokenRepo{}, &mockPasswordHasher{})

		err := svc.Register(ctx, RegisterRequest{Email: "alice@example.com", Password: ""})
		if !errors.Is(err, ErrInvalidPassword) {
			t.Fatalf("expected ErrInvalidPassword, got %v", err)
		}
	})

	t.Run("email already exists returns ErrEmailExists", func(t *testing.T) {
		userRepo := &mockUserRepo{
			existsByEmailFn: func(_ context.Context, _ string) (bool, error) { return true, nil },
		}
		svc := newTestService(userRepo, &mockTokenRepo{}, &mockPasswordHasher{})

		err := svc.Register(ctx, RegisterRequest{Email: "alice@example.com", Password: "secret"})
		if !errors.Is(err, ErrEmailExists) {
			t.Fatalf("expected ErrEmailExists, got %v", err)
		}
	})

	t.Run("email is normalised to lowercase", func(t *testing.T) {
		var capturedEmail string
		userRepo := &mockUserRepo{
			existsByEmailFn: func(_ context.Context, email string) (bool, error) {
				capturedEmail = email
				return false, nil
			},
			createFn: func(_ context.Context, _ User) error { return nil },
		}
		hasher := &mockPasswordHasher{
			hashFn: func(_ string) (string, error) { return "hashed", nil },
		}
		svc := newTestService(userRepo, &mockTokenRepo{}, hasher)

		_ = svc.Register(ctx, RegisterRequest{Email: "  Alice@Example.COM  ", Password: "secret"})
		if capturedEmail != "alice@example.com" {
			t.Fatalf("expected normalised email, got %q", capturedEmail)
		}
	})

	t.Run("repository ExistsByEmail error is propagated", func(t *testing.T) {
		dbErr := errors.New("db connection error")
		userRepo := &mockUserRepo{
			existsByEmailFn: func(_ context.Context, _ string) (bool, error) { return false, dbErr },
		}
		svc := newTestService(userRepo, &mockTokenRepo{}, &mockPasswordHasher{})

		err := svc.Register(ctx, RegisterRequest{Email: "alice@example.com", Password: "secret"})
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected db error, got %v", err)
		}
	})

	t.Run("password hasher error is propagated", func(t *testing.T) {
		hashErr := errors.New("bcrypt failure")
		userRepo := &mockUserRepo{
			existsByEmailFn: func(_ context.Context, _ string) (bool, error) { return false, nil },
		}
		hasher := &mockPasswordHasher{
			hashFn: func(_ string) (string, error) { return "", hashErr },
		}
		svc := newTestService(userRepo, &mockTokenRepo{}, hasher)

		err := svc.Register(ctx, RegisterRequest{Email: "alice@example.com", Password: "secret"})
		if !errors.Is(err, hashErr) {
			t.Fatalf("expected hash error, got %v", err)
		}
	})

	t.Run("repository Create error is propagated", func(t *testing.T) {
		createErr := errors.New("insert failed")
		userRepo := &mockUserRepo{
			existsByEmailFn: func(_ context.Context, _ string) (bool, error) { return false, nil },
			createFn:        func(_ context.Context, _ User) error { return createErr },
		}
		hasher := &mockPasswordHasher{
			hashFn: func(_ string) (string, error) { return "hashed", nil },
		}
		svc := newTestService(userRepo, &mockTokenRepo{}, hasher)

		err := svc.Register(ctx, RegisterRequest{Email: "alice@example.com", Password: "secret"})
		if !errors.Is(err, createErr) {
			t.Fatalf("expected create error, got %v", err)
		}
	})
}

// ---------------------------------------------------------------------------
// Login tests
// ---------------------------------------------------------------------------

func TestService_Login(t *testing.T) {
	ctx := context.Background()

	baseUser := User{ID: 1, Name: "Alice", Email: "alice@example.com", HashPassword: "hashed"}

	t.Run("success", func(t *testing.T) {
		userRepo := &mockUserRepo{
			getByEmailFn: func(_ context.Context, _ string) (User, error) { return baseUser, nil },
		}
		tokenRepo := &mockTokenRepo{
			storeRefreshTokenFn: func(_ context.Context, _ int64, _ string, _ time.Time) error { return nil },
		}
		hasher := &mockPasswordHasher{
			compareFn: func(_, _ string) error { return nil },
		}
		svc := newTestService(userRepo, tokenRepo, hasher)

		user, pair, err := svc.Login(ctx, LoginRequest{Email: "alice@example.com", Password: "secret"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user.Email != baseUser.Email {
			t.Errorf("expected email %q, got %q", baseUser.Email, user.Email)
		}
		if pair.AccessToken == "" || pair.RefreshToken == "" {
			t.Error("expected non-empty token pair")
		}
	})

	t.Run("empty email returns ErrInvalidEmail", func(t *testing.T) {
		svc := newTestService(&mockUserRepo{}, &mockTokenRepo{}, &mockPasswordHasher{})

		_, _, err := svc.Login(ctx, LoginRequest{Email: "", Password: "secret"})
		if !errors.Is(err, ErrInvalidEmail) {
			t.Fatalf("expected ErrInvalidEmail, got %v", err)
		}
	})

	t.Run("empty password returns ErrInvalidPassword", func(t *testing.T) {
		svc := newTestService(&mockUserRepo{}, &mockTokenRepo{}, &mockPasswordHasher{})

		_, _, err := svc.Login(ctx, LoginRequest{Email: "alice@example.com", Password: ""})
		if !errors.Is(err, ErrInvalidPassword) {
			t.Fatalf("expected ErrInvalidPassword, got %v", err)
		}
	})

	t.Run("user not found returns ErrUserNotFound", func(t *testing.T) {
		userRepo := &mockUserRepo{
			getByEmailFn: func(_ context.Context, _ string) (User, error) {
				return User{}, sql.ErrNoRows
			},
		}
		svc := newTestService(userRepo, &mockTokenRepo{}, &mockPasswordHasher{})

		_, _, err := svc.Login(ctx, LoginRequest{Email: "unknown@example.com", Password: "secret"})
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("wrong password returns ErrUserNotFound", func(t *testing.T) {
		userRepo := &mockUserRepo{
			getByEmailFn: func(_ context.Context, _ string) (User, error) { return baseUser, nil },
		}
		hasher := &mockPasswordHasher{
			compareFn: func(_, _ string) error { return errors.New("mismatch") },
		}
		svc := newTestService(userRepo, &mockTokenRepo{}, hasher)

		_, _, err := svc.Login(ctx, LoginRequest{Email: "alice@example.com", Password: "wrong"})
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("repository GetByEmail error is propagated", func(t *testing.T) {
		dbErr := errors.New("db error")
		userRepo := &mockUserRepo{
			getByEmailFn: func(_ context.Context, _ string) (User, error) { return User{}, dbErr },
		}
		svc := newTestService(userRepo, &mockTokenRepo{}, &mockPasswordHasher{})

		_, _, err := svc.Login(ctx, LoginRequest{Email: "alice@example.com", Password: "secret"})
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected db error, got %v", err)
		}
	})

	t.Run("StoreRefreshToken error is propagated", func(t *testing.T) {
		storeErr := errors.New("store failed")
		userRepo := &mockUserRepo{
			getByEmailFn: func(_ context.Context, _ string) (User, error) { return baseUser, nil },
		}
		tokenRepo := &mockTokenRepo{
			storeRefreshTokenFn: func(_ context.Context, _ int64, _ string, _ time.Time) error { return storeErr },
		}
		hasher := &mockPasswordHasher{
			compareFn: func(_, _ string) error { return nil },
		}
		svc := newTestService(userRepo, tokenRepo, hasher)

		_, _, err := svc.Login(ctx, LoginRequest{Email: "alice@example.com", Password: "secret"})
		if !errors.Is(err, storeErr) {
			t.Fatalf("expected store error, got %v", err)
		}
	})
}

// ---------------------------------------------------------------------------
// Refresh tests
// ---------------------------------------------------------------------------

func TestService_Refresh(t *testing.T) {
	ctx := context.Background()

	// Generate a real refresh token to use across subtests.
	validPair, err := GenerateTokenPair(1, "alice@example.com", "access-secret", "refresh-secret", 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("setup: failed to generate token pair: %v", err)
	}
	validHash := HashToken(validPair.RefreshToken)

	t.Run("success rotates tokens", func(t *testing.T) {
		storedToken := RefreshToken{
			ID:        1,
			UserID:    1,
			TokenHash: validHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		}

		var deletedHash string
		var storedHash string

		tokenRepo := &mockTokenRepo{
			getRefreshTokenFn: func(_ context.Context, hash string) (RefreshToken, error) {
				return storedToken, nil
			},
			deleteRefreshTokenFn: func(_ context.Context, hash string) error {
				deletedHash = hash
				return nil
			},
			storeRefreshTokenFn: func(_ context.Context, _ int64, hash string, _ time.Time) error {
				storedHash = hash
				return nil
			},
		}
		svc := newTestService(&mockUserRepo{}, tokenRepo, &mockPasswordHasher{})

		pair, err := svc.Refresh(ctx, validPair.RefreshToken)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if pair.AccessToken == "" || pair.RefreshToken == "" {
			t.Error("expected non-empty new token pair")
		}
		if deletedHash != validHash {
			t.Errorf("expected old token hash to be deleted, got %q", deletedHash)
		}
		if storedHash == "" {
			t.Error("expected a new token hash to be stored")
		}
	})

	t.Run("invalid token string returns ErrInvalidToken", func(t *testing.T) {
		svc := newTestService(&mockUserRepo{}, &mockTokenRepo{}, &mockPasswordHasher{})

		_, err := svc.Refresh(ctx, "not.a.valid.jwt")
		if !errors.Is(err, ErrInvalidToken) {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("token not in DB returns ErrInvalidToken", func(t *testing.T) {
		tokenRepo := &mockTokenRepo{
			getRefreshTokenFn: func(_ context.Context, _ string) (RefreshToken, error) {
				return RefreshToken{}, sql.ErrNoRows
			},
		}
		svc := newTestService(&mockUserRepo{}, tokenRepo, &mockPasswordHasher{})

		_, err := svc.Refresh(ctx, validPair.RefreshToken)
		if !errors.Is(err, ErrInvalidToken) {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("expired stored token returns ErrInvalidToken", func(t *testing.T) {
		expiredToken := RefreshToken{
			TokenHash: validHash,
			ExpiresAt: time.Now().Add(-time.Hour), // already expired
		}
		tokenRepo := &mockTokenRepo{
			getRefreshTokenFn: func(_ context.Context, _ string) (RefreshToken, error) {
				return expiredToken, nil
			},
			deleteRefreshTokenFn: func(_ context.Context, _ string) error { return nil },
		}
		svc := newTestService(&mockUserRepo{}, tokenRepo, &mockPasswordHasher{})

		_, err := svc.Refresh(ctx, validPair.RefreshToken)
		if !errors.Is(err, ErrInvalidToken) {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("GetRefreshToken db error is propagated", func(t *testing.T) {
		dbErr := errors.New("db error")
		tokenRepo := &mockTokenRepo{
			getRefreshTokenFn: func(_ context.Context, _ string) (RefreshToken, error) {
				return RefreshToken{}, dbErr
			},
		}
		svc := newTestService(&mockUserRepo{}, tokenRepo, &mockPasswordHasher{})

		_, err := svc.Refresh(ctx, validPair.RefreshToken)
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected db error, got %v", err)
		}
	})

	t.Run("DeleteRefreshToken error is propagated", func(t *testing.T) {
		deleteErr := errors.New("delete failed")
		storedToken := RefreshToken{
			TokenHash: validHash,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		tokenRepo := &mockTokenRepo{
			getRefreshTokenFn: func(_ context.Context, _ string) (RefreshToken, error) {
				return storedToken, nil
			},
			deleteRefreshTokenFn: func(_ context.Context, _ string) error { return deleteErr },
		}
		svc := newTestService(&mockUserRepo{}, tokenRepo, &mockPasswordHasher{})

		_, err := svc.Refresh(ctx, validPair.RefreshToken)
		if !errors.Is(err, deleteErr) {
			t.Fatalf("expected delete error, got %v", err)
		}
	})
}

// ---------------------------------------------------------------------------
// Logout tests
// ---------------------------------------------------------------------------

func TestService_Logout(t *testing.T) {
	ctx := context.Background()

	t.Run("success deletes token", func(t *testing.T) {
		var deletedHash string
		tokenRepo := &mockTokenRepo{
			deleteRefreshTokenFn: func(_ context.Context, hash string) error {
				deletedHash = hash
				return nil
			},
		}
		svc := newTestService(&mockUserRepo{}, tokenRepo, &mockPasswordHasher{})

		rawToken := "some-refresh-token"
		err := svc.Logout(ctx, rawToken)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if deletedHash != HashToken(rawToken) {
			t.Errorf("expected token hash %q to be deleted, got %q", HashToken(rawToken), deletedHash)
		}
	})

	t.Run("token already gone is treated as success", func(t *testing.T) {
		tokenRepo := &mockTokenRepo{
			deleteRefreshTokenFn: func(_ context.Context, _ string) error { return sql.ErrNoRows },
		}
		svc := newTestService(&mockUserRepo{}, tokenRepo, &mockPasswordHasher{})

		err := svc.Logout(ctx, "any-token")
		if err != nil {
			t.Fatalf("expected no error when token already deleted, got %v", err)
		}
	})

	t.Run("delete db error is propagated", func(t *testing.T) {
		deleteErr := errors.New("db failure")
		tokenRepo := &mockTokenRepo{
			deleteRefreshTokenFn: func(_ context.Context, _ string) error { return deleteErr },
		}
		svc := newTestService(&mockUserRepo{}, tokenRepo, &mockPasswordHasher{})

		err := svc.Logout(ctx, "any-token")
		if !errors.Is(err, deleteErr) {
			t.Fatalf("expected delete error, got %v", err)
		}
	})
}
