package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lz1marine/notification-service/pkg/channels"
	"github.com/lz1marine/notification-service/pkg/controllers"
	"github.com/lz1marine/notification-service/pkg/queue"
)

func main() {
	channelHandler := getChannelHandler()
	channelHandler.Start()

	// TODO: handle graceful shutdown
}

func getChannelHandler() *controllers.ChannelHandler {
	username, password, host, port, maxConnections, qEp, qPwd, err := readCredentials()
	if err != nil {
		panic(err)
	}

	// TODO: log info
	fmt.Printf("username: %s\nhost: %s\nport: %d\n", username, host, port)
	channel := channels.NewEmailChannel(host, port, username, password, uint(maxConnections))

	redis := queue.NewRedis(qEp, qPwd, 0, 5)

	return controllers.NewChannelHandler(channel, redis)
}

func readCredentials() (string, string, string, int, int, string, string, error) {
	// Read username from file
	usernameBytes, err := os.ReadFile("/home/marin/dev/secrets/username")
	if err != nil {
		return "", "", "", 0, 0, "", "", fmt.Errorf("failed to read username file: %v", err)
	}
	username := string(usernameBytes)

	// Read password from file
	passwordBytes, err := os.ReadFile("/home/marin/dev/secrets/password")
	if err != nil {
		return "", "", "", 0, 0, "", "", fmt.Errorf("failed to read password file: %v", err)
	}
	password := string(passwordBytes)

	// Read password from file
	qPasswordBytes, err := os.ReadFile("/home/marin/dev/secrets/queue_password")
	if err != nil {
		return "", "", "", 0, 0, "", "", fmt.Errorf("failed to read the queue password file: %v", err)
	}
	qPassword := string(qPasswordBytes)

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
		return "", "", "", 0, 0, "", "", fmt.Errorf("failed to adapt port to a number: %v", err)
	}

	// Read max connections
	maxConnectionsString := os.Getenv("MAX_CONNECTIONS")
	if maxConnectionsString == "" {
		maxConnectionsString = "200"
	}
	maxConnections, err := strconv.Atoi(maxConnectionsString)
	if err != nil {
		return "", "", "", 0, 0, "", "", fmt.Errorf("failed to adapt max connections to a number: %v", err)
	}

	// Read distributed queue location from file
	qEndpoint := os.Getenv("QUEUE_ENDPOINT")
	if qEndpoint == "" {
		qEndpoint = "localhost:6379"
	}

	return username, password, host, port, maxConnections, qEndpoint, qPassword, nil
}
