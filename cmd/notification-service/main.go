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

// @title           Notification Server API
// @version         1.0
// @description     This is the notification server API that handles both external user notification subscriptions and internal notifications
// @termsOfService  http://swagger.io/terms/

// @license.name  GNU General Public License v3.0
// @license.url   https://github.com/lz1marine/notification-service/?tab=GPL-3.0-1-ov-file

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BearerAuth
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
	err = engine.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}

	// TODO: handle graceful shutdown
}

func prepareEngine() *gin.Engine {
	return gin.Default()
}

func readConfig() (string, string, string, string, string, int, error) {
	// Read password from file
	qPasswordBytes, err := os.ReadFile("/app/secrets/queue_password")
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
	dbUsernameBytes, err := os.ReadFile("/app/secrets/database_username")
	if err != nil {
		dbUsername = "admin"
	} else {
		dbUsername = string(dbUsernameBytes)
	}

	// Read the db password from file
	dbPassword := ""
	dbPasswordBytes, err := os.ReadFile("/app/secrets/database_password")
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
