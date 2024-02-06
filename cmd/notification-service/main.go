package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lz1marine/notification-service/pkg/client"
	dbctrl "github.com/lz1marine/notification-service/pkg/database/controller"
	eh "github.com/lz1marine/notification-service/pkg/handler"
	ih "github.com/lz1marine/notification-service/pkg/handler/inter"
	"github.com/lz1marine/notification-service/pkg/http"
	"github.com/lz1marine/notification-service/pkg/queue"
)

func main() {
	// Create a new Gin router
	engine := prepareEngine()

	// Read the config
	qEp, qPwd, dbUsername, dbPassword, dbLocation, dbPort, err := readConfig()
	if err != nil {
		panic(err)
	}

	db := client.NewDBClient(dbUsername, dbPassword, dbLocation, dbPort)
	dbChanController := dbctrl.NewChannelController(db)
	dbMessageController := dbctrl.NewMessageController(db)
	dbUserController := dbctrl.NewUserController(db)

	// Setup routes
	http.NotificationServiceHandlers(engine, eh.NewExternalHandler(dbChanController))

	redisQ := queue.NewRedis(qEp, qPwd, 0, 5)
	redisB := client.NewRedisBackup(qEp, qPwd, 5)

	http.NotificationHandlers(engine, ih.NewNotificationHandler(redisQ, redisB, dbChanController, dbMessageController, dbUserController))

	// Start the server
	engine.Run("localhost:8080")

	// TODO: handle graceful shutdown
}

func prepareEngine() *gin.Engine {
	return gin.Default()
}

func readConfig() (string, string, string, string, string, int, error) {
	// Read password from file
	qPasswordBytes, err := os.ReadFile("/home/marin/dev/secrets/queue_password")
	if err != nil {
		return "", "", "", "", "", 0, fmt.Errorf("failed to read the queue password file: %v", err)
	}
	qPassword := string(qPasswordBytes)

	// Read distributed queue location from file
	qEndpoint := os.Getenv("QUEUE_ENDPOINT")
	if qEndpoint == "" {
		qEndpoint = "localhost:6379"
	}

	// Read the db username from file
	dbUsername := ""
	dbUsernameBytes, err := os.ReadFile("/home/marin/dev/secrets/database_username")
	if err != nil {
		dbUsername = "admin"
	} else {
		dbUsername = string(dbUsernameBytes)
	}

	// Read the db password from file
	dbPassword := ""
	dbPasswordBytes, err := os.ReadFile("/home/marin/dev/secrets/database_password")
	if err != nil {
		dbPassword = "admin"
	} else {
		dbPassword = string(dbPasswordBytes)
	}

	// Read distributed queue location from file
	dbLocation := os.Getenv("DB_LOCATION")
	if dbLocation == "" {
		dbLocation = "localhost"
	}

	// Read distributed queue location from file
	dbPortString := os.Getenv("DB_PORT")
	if dbPortString == "" {
		dbPortString = "3306"
	}
	dbPort, err := strconv.Atoi(dbPortString)
	if err != nil {
		return "", "", "", "", "", 0, fmt.Errorf("failed to adapt port to a number: %v", err)
	}

	return qEndpoint, qPassword, dbUsername, dbPassword, dbLocation, dbPort, nil
}
