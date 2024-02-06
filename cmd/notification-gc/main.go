package main

import (
	"fmt"
	"os"

	"github.com/lz1marine/notification-service/pkg/client"
	"github.com/lz1marine/notification-service/pkg/controller"
	"github.com/lz1marine/notification-service/pkg/queue"
)

func main() {
	gc := getGarbageCollector()
	gc.Start()

	// TODO: handle graceful shutdown
}

func getGarbageCollector() *controller.GarbageCollector {
	qEp, qPwd, err := readCredentials()
	if err != nil {
		panic(err)
	}

	// TODO: get databases from config, not 0, 5, 10
	redisQ := queue.NewRedis(qEp, qPwd, 0, 5)
	redisB := client.NewRedisBackup(qEp, qPwd, 5)

	fmt.Println("starting garbage collector")
	return controller.NewGarbageCollector(redisB, redisQ)
}

func readCredentials() (string, string, error) {
	// Read password from file
	qPasswordBytes, err := os.ReadFile("/app/secrets/queue_password")
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
