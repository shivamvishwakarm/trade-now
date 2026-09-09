package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	accessTokenCookie  = "access_token"
	refreshTokenCookie = "refresh_token"
)

type Handler struct {
	Logger  *zap.Logger
	Service *Service
}

type HandlerDeps struct {
	Logger  *zap.Logger
	Service *Service
}

func NewHandler(deps HandlerDeps) *Handler {
	return &Handler{
		Logger:  deps.Logger,
		Service: deps.Service,
	}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user account. The email is normalised (lowercased + trimmed) before storage. Passwords are hashed with bcrypt.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		RegisterRequest	true	"Registration payload"
//	@Success		201		{object}	RegisterResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body or missing required fields"
//	@Failure		409		{object}	ErrorResponse	"A user with that email already exists"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/auth/register [post]
func (h *Handler) Register(ctx *gin.Context) {
	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn("invalid registration request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.Service.Register(ctx, req); err != nil {
		switch err {
		case ErrInvalidEmail, ErrInvalidPassword:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case ErrEmailExists:
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, RegisterResponse{Message: "user registered successfully"})
}

// Login godoc
//
//	@Summary		Log in
//	@Description	Authenticates a user with email and password. On success, sets two HttpOnly cookies: `access_token` (path `/`) and `refresh_token` (path `/auth/refresh`). The raw tokens are never returned in the response body.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginRequest	true	"Login credentials"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body"
//	@Failure		401		{object}	ErrorResponse	"Invalid email or password"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/auth/login [post]
func (h *Handler) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn("invalid login request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, pair, err := h.Service.Login(ctx, req)
	if err != nil {
		switch err {
		case ErrInvalidEmail, ErrInvalidPassword, ErrUserNotFound:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	setTokenCookies(ctx, pair, h.Service.accessExpiry, h.Service.refreshExpiry)

	ctx.JSON(http.StatusOK, LoginResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

// Refresh godoc
//
//	@Summary		Refresh access token
//	@Description	Rotates the token pair. Reads the `refresh_token` HttpOnly cookie, validates it, deletes the old token from the database, and issues a fresh access + refresh token pair (stored as cookies). The old refresh token is invalidated immediately.
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	MessageResponse	"Tokens refreshed successfully"
//	@Failure		401	{object}	ErrorResponse	"Refresh token missing, invalid, or expired"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Security		RefreshTokenCookie
//	@Router			/auth/refresh [post]
func (h *Handler) Refresh(ctx *gin.Context) {
	rawToken, err := ctx.Cookie(refreshTokenCookie)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	pair, err := h.Service.Refresh(ctx, rawToken)
	if err != nil {
		clearTokenCookies(ctx)
		switch err {
		case ErrInvalidToken:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	setTokenCookies(ctx, pair, h.Service.accessExpiry, h.Service.refreshExpiry)
	ctx.JSON(http.StatusOK, gin.H{"message": "tokens refreshed"})
}

// Logout godoc
//
//	@Summary		Log out
//	@Description	Revokes the current refresh token (best-effort DB deletion) and clears both `access_token` and `refresh_token` cookies by setting their MaxAge to -1. Always returns 200 — even if no cookie was present.
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	MessageResponse	"Logged out successfully"
//	@Security		RefreshTokenCookie
//	@Router			/auth/logout [post]
func (h *Handler) Logout(ctx *gin.Context) {
	rawToken, err := ctx.Cookie(refreshTokenCookie)
	if err == nil {
		// Best-effort revocation — ignore error so the cookie is always cleared
		_ = h.Service.Logout(ctx, rawToken)
	}

	clearTokenCookies(ctx)
	ctx.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// --- helpers ---

func setTokenCookies(ctx *gin.Context, pair TokenPair, accessExpiry, refreshExpiry time.Duration) {
	ctx.SetCookie(
		accessTokenCookie,
		pair.AccessToken,
		int(accessExpiry.Seconds()),
		"/",
		"",
		false, // set to true in prod (HTTPS only)
		true,  // HttpOnly
	)
	ctx.SetCookie(
		refreshTokenCookie,
		pair.RefreshToken,
		int(refreshExpiry.Seconds()),
		"/auth/refresh",
		"",
		false,
		true,
	)
}

func clearTokenCookies(ctx *gin.Context) {
	ctx.SetCookie(accessTokenCookie, "", -1, "/", "", false, true)
	ctx.SetCookie(refreshTokenCookie, "", -1, "/auth/refresh", "", false, true)
}
