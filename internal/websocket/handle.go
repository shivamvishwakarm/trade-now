package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Handler struct {
	logger   *zap.Logger
	upgrader *websocket.Upgrader
}

type HandlerDeps struct {
	Logger *zap.Logger
}

func NewHandler(deps HandlerDeps) *Handler {
	return &Handler{
		logger: deps.Logger,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (h *Handler) Handle(ctx *gin.Context) {

	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		h.logger.Error(
			"websocket upgrade fail",
			zap.Error(err),
		)
		return
	}
	defer conn.Close()

	messageType, message, err := conn.ReadMessage()

	if err != nil {
		h.logger.Error(
			"websocket read fail",
			zap.Error(err),
		)
		return
	}

	h.logger.Info(
		"websocket message receive",
		zap.Int("message_type", messageType),
		zap.ByteString("message", message),
	)

	if err := conn.WriteMessage(
		websocket.TextMessage,
		[]byte(`{"message":"connected to websocket"}`),
	); err != nil {
		h.logger.Error(
			"websocket write failed",
			zap.Error(err),
		)
		return
	}

}
