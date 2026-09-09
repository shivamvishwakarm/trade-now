package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func (h *Handler) Register(ctx *gin.Context) {
	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn(
			"invalid registraction request",
			zap.Error(err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	if err := h.Service.Register(ctx, req); err != nil {
		// err maping
		return
	}

	ctx.JSON(http.StatusCreated, RegisterResponse{
		Message: "user registered successfully",
	})
}

func (h *Handler) Login(ctx *gin.Context) {

	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn(
			"invalid login request",
			zap.Error(err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	user, err := h.Service.Login(ctx, req)

	if err != nil {
		//
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		Email: user.Email,
		Name:  user.Name,
	})
}

func (h *Handler) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "logout not implemented"})
}

func (h *Handler) Refresh(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "refresh not implemented"})
}
