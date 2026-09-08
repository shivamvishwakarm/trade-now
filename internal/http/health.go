package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func registerHealthRoutes(router *gin.Engine) {

	router.GET("/health", healthHandler)
}

func healthHandler(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
	})
}
