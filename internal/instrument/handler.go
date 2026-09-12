package instrument

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

func (i *Handler) GetBySymbol(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "not implemented, gus",
	})
}
