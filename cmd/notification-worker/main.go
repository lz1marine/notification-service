package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lz1marine/notification-service/pkg/channels"
	handlers "github.com/lz1marine/notification-service/pkg/handlers/temp"
	"github.com/lz1marine/notification-service/pkg/http"
)

func main() {
	// Create a new Gin router
	engine := prepareEngine()

	channelHandler := getChannelHandler()

	// Setup routes
	http.TempNotificationWorkerHandlers(engine, channelHandler)

	// Start the server
	engine.Run("localhost:8080")
}

func prepareEngine() *gin.Engine {
	return gin.Default()
}

func getChannelHandler() *handlers.ChannelHandler {
	username, password, host, port, maxConnections, err := readCredentials()
	if err != nil {
		panic(err)
	}

	// TODO: log info
	fmt.Printf("username: %s\npassword: %s\nhost: %s\nport: %d\n", username, password, host, port)
	return handlers.NewChannelHandler(channels.NewEmailChannel(host, port, username, password, uint(maxConnections)))
}

func readCredentials() (string, string, string, int, int, error) {
	// Read username from file
	usernameBytes, err := os.ReadFile("/home/marin/dev/secrets/username")
	if err != nil {
		return "", "", "", 0, 0, fmt.Errorf("failed to read username file: %v", err)
	}
	username := string(usernameBytes)

	// Read password from file
	passwordBytes, err := os.ReadFile("/home/marin/dev/secrets/password")
	if err != nil {
		return "", "", "", 0, 0, fmt.Errorf("failed to read password file: %v", err)
	}
	password := string(passwordBytes)

	// Read host from file
	host := os.Getenv("HOST")
	if host == "" {
		host = "smtp.gmail.com"
	}

	// Read port
	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "587"
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		return "", "", "", 0, 0, fmt.Errorf("failed to adapt port to a number: %v", err)
	}

	// Read max connections
	maxConnectionsString := os.Getenv("MAX_CONNECTIONS")
	if maxConnectionsString == "" {
		maxConnectionsString = "200"
	}
	maxConnections, err := strconv.Atoi(maxConnectionsString)
	if err != nil {
		return "", "", "", 0, 0, fmt.Errorf("failed to adapt max connections to a number: %v", err)
	}

	return username, password, host, port, maxConnections, nil
}
