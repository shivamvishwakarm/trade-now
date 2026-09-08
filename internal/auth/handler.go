package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	Logger *zap.Logger
}

type HandlerDeps struct {
	Logger         *zap.Logger
	PasswordHasher PasswordHasher
}

func NewHandler(deps HandlerDeps) *Handler {
	return &Handler{
		Logger: deps.Logger,
	}
}

func (H *Handler) Register(ctx *gin.Context) {

	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		H.Logger.Warn(
			"invalid registeration request",
			zap.Error(err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "register not implemented",
	})
}

func (H *Handler) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "login not implemnted",
	})
}

func (H *Handler) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "logout not implemented",
	})
}

func (H *Handler) Refresh(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "refresh not implemented",
	})
}
