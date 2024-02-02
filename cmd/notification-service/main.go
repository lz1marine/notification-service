package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lz1marine/notification-service/pkg/http"
)

func main() {
	// Create a new Gin router
	engine := prepareEngine()

	// Setup routes
	http.AddHandlers(engine)

	// Start the server
	engine.Run("localhost:8080")
}

func prepareEngine() *gin.Engine {
	return gin.Default()
}
