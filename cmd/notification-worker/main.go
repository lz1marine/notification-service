package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lz1marine/notification-service/pkg/channel"
	"github.com/lz1marine/notification-service/pkg/client"
	"github.com/lz1marine/notification-service/pkg/controller"
	"github.com/lz1marine/notification-service/pkg/queue"
)

func main() {
	channelHandler := getChannelHandler()
	channelHandler.Start()

	// TODO: handle graceful shutdown
}

func getChannelHandler() *controller.ChannelHandler {
	// Read type of channel
	chanType := os.Getenv("TYPE")
	if chanType == "" {
		panic("TYPE not set")
	}

	maxConnections, qEp, qPwd, err := readCredentials(chanType)
	if err != nil {
		panic(err)
	}

	var ch channel.Channel
	switch chanType {
	case "email":
		ch = getNewEmailChannel()
	case "sms":
		ch = getNewSMSChannel()
	default:
		panic(fmt.Sprintf("unknown channel type: %s", chanType))
	}

	// TODO: get databases from config, not 0, 5, 10
	redisQ := queue.NewRedis(qEp, qPwd, 0, 5)
	redisB := client.NewRedisBackup(qEp, qPwd, 5)
	redisT := client.NewRedisTemplate(qEp, qPwd, 10)

	return controller.NewChannelHandler(ch, redisQ, redisT, redisB, uint(maxConnections))
}

func readCredentials(chanType string) (int, string, string, error) {
	// Read password from file
	qPasswordBytes, err := os.ReadFile("/app/secrets/queue_password")
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to read the queue password file: %v", err)
	}
	qPassword := string(qPasswordBytes)

	// Read max connections
	maxConnectionsString := os.Getenv("MAX_CONNECTIONS")
	if maxConnectionsString == "" {
		maxConnectionsString = "200"
	}
	maxConnections, err := strconv.Atoi(maxConnectionsString)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to adapt max connections to a number: %v", err)
	}

	// Read distributed queue location from file
	qEndpoint := os.Getenv("QUEUE_ENDPOINT")
	if qEndpoint == "" {
		qEndpoint = "localhost:6379"
	}

	return maxConnections, qEndpoint, qPassword, nil
}

func readEmailCredentials() (string, string, string, int, error) {
	// Read username from file
	usernameBytes, err := os.ReadFile("/app/secrets/email_username")
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to read username file: %v", err)
	}
	username := string(usernameBytes)

	// Read password from file
	passwordBytes, err := os.ReadFile("/app/secrets/email_password")
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to read password file: %v", err)
	}
	password := string(passwordBytes)

	// Read host
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
		return "", "", "", 0, fmt.Errorf("failed to adapt port to a number: %v", err)
	}

	return username, password, host, port, nil
}

func readSMSCredentials() (string, string, string, error) {
	// Read username from file
	usernameBytes, err := os.ReadFile("/app/secrets/sms_username")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read username file: %v", err)
	}
	username := string(usernameBytes)

	// Read password from file
	passwordBytes, err := os.ReadFile("/app/secrets/sms_password")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read password file: %v", err)
	}
	password := string(passwordBytes)

	// Read password from file
	senderBytes, err := os.ReadFile("/app/secrets/sms_phone_sender")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read password file: %v", err)
	}
	sender := string(senderBytes)

	return username, password, sender, nil
}

func getNewEmailChannel() channel.Channel {
	username, password, host, port, err := readEmailCredentials()
	if err != nil {
		panic(err)
	}

	fmt.Printf("username: %s\nhost: %s\nport: %d\n", username, host, port)

	return channel.NewEmailChannel(host, port, username, password)
}

func getNewSMSChannel() channel.Channel {
	username, password, sender, err := readSMSCredentials()
	if err != nil {
		panic(err)
	}

	fmt.Printf("username: %s\nsender: %s\n", username, sender)

	return channel.NewSMSChannel(username, password, sender)
}
