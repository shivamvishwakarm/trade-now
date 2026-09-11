package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *zap.Logger
	service *Service
}

type HandlerDeps struct {
	Logger  *zap.Logger
	Service *Service
}

func NewHandler(deps HandlerDeps) *Handler {
	return &Handler{
		logger:  deps.Logger,
		service: deps.Service,
	}
}

// GetMe godoc
//
//	@Summary		Get current user profile
//	@Description	Returns the profile of the authenticated user.
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	Profile
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Security		AccessTokenCookie
//	@Router			/api/v1/me [get]
//
// TODO: implement
func (h *Handler) GetMe(ctx *gin.Context) {

	user, err := h.service.GetProfile(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized access",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
