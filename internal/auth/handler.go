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
	Logger *zap.Logger
}

func NewHandler(deps HandlerDeps) *Handler {
	return &Handler{
		Logger: deps.Logger,
	}
}

func (H *Handler) Register(ctx *gin.Context) {

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
