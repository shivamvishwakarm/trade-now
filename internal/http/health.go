package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func registerHealthRoutes(router *gin.Engine) {
	router.GET("/health", healthHandler)
}

// healthHandler godoc
//
//	@Summary		Health check
//	@Description	Returns the current server status and UTC timestamp. Use this to confirm the API is reachable.
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"status ok"
//	@Router			/health [get]
func healthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
	})
}
