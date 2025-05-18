package app

import (
	"orbital/internal/handler" // Import the handler package

	"github.com/gin-gonic/gin"
)

// NewRouter creates and returns a new Gin router
func NewRouter() *gin.Engine {
	r := gin.Default()

	// Register the /createArgoApp route
	r.POST("/upsertArgoApp", handler.UpsertArgoApp)

	return r
}
