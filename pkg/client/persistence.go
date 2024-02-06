package client

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBClient struct {
	notifications *gorm.DB
	users         *gorm.DB
}

// NewDBClient creates a new client and connects to the database. This connection is reusable but it has to be closed.
// TODO: use config.DB instead, we can also add a connection pool (SetMaxIdleConns, SetMaxOpenConns, SetConnMaxLifetime)
func NewDBClient(username, password, location string, port int) *DBClient {
	dsnNotifications := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, location, port, "notifications")
	dsnUsers := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, location, port, "users")

	// Open the database connections
	dbNotifications, err := gorm.Open(mysql.Open(dsnNotifications), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to database %s: %v", "notifications", err)
	}
	dbUsers, err := gorm.Open(mysql.Open(dsnUsers), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to database %s: %v", "users", err)
	}

	return &DBClient{
		notifications: dbNotifications,
		users:         dbUsers,
	}
}

// TODO: code duplication
// Close closes the database connection.
func (dbc *DBClient) Close() error {
	// Close the database connections
	err := close(dbc.notifications)
	if err != nil {
		return err
	}

	err = close(dbc.users)
	if err != nil {
		return err
	}

	return nil
}

func (dbc *DBClient) GetNotificationsDB() *gorm.DB {
	return dbc.notifications
}

func (dbc *DBClient) GetUsersDB() *gorm.DB {
	return dbc.users
}

func close(db *gorm.DB) error {
	// Close the database connection
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Failed to get DB instance: %v", err)
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		fmt.Printf("Failed to close DB connection: %v", err)
		return err
	}

	return nil
}
