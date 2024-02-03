package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	ih "github.com/lz1marine/notification-service/pkg/handlers/inter"
	"github.com/lz1marine/notification-service/pkg/http"
	"github.com/lz1marine/notification-service/pkg/queue"
)

func main() {
	// Create a new Gin router
	engine := prepareEngine()

	// Setup routes
	http.NotificationServiceHandlers(engine)

	// TODO: should be internal
	qEp, qPwd, err := readConfig()
	if err != nil {
		panic(err)
	}
	redis := queue.NewRedis(qEp, qPwd, 0, 5)
	http.NotificationHandlers(engine, ih.NewNotificationHandler(redis))

	// Start the server
	engine.Run("localhost:8080")

	// TODO: handle graceful shutdown
}

func prepareEngine() *gin.Engine {
	return gin.Default()
}

func readConfig() (string, string, error) {
	// Read password from file
	qPasswordBytes, err := os.ReadFile("/home/marin/dev/secrets/queue_password")
	if err != nil {
		return "", "", fmt.Errorf("failed to read the queue password file: %v", err)
	}
	qPassword := string(qPasswordBytes)

	// Read distributed queue location from file
	qEndpoint := os.Getenv("QUEUE_ENDPOINT")
	if qEndpoint == "" {
		qEndpoint = "localhost:6379"
	}

	return qEndpoint, qPassword, nil
}
