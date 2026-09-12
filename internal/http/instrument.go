package http

import "github.com/gin-gonic/gin"

func registerInstrucmentRoute(router *gin.RouterGroup, deps RouterDeps) {

	instruments := router.Group("/instruments")

	// instruments.GET("", deps.InstrumentHandler.List)
	instruments.GET("/:symbol", deps.InstrumentHandler.GetBySymbol)
}
