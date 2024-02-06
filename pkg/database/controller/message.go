package controller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lz1marine/notification-service/pkg/client"
	"github.com/lz1marine/notification-service/pkg/database/entity"
)

// MessageController is the controller for messages
type MessageController struct {
	db *client.DBClient
}

// NewMessageController creates a new MessageController
func NewMessageController(db *client.DBClient) *MessageController {
	return &MessageController{
		db: db,
	}
}

// SetMessageStatus updates the state of a message
func (mc *MessageController) SetMessageStatus(ctx context.Context, messageID string, status int) error {
	var message entity.Message
	result := mc.db.GetNotificationsDB().WithContext(ctx).Where("id = ?", messageID).Find(&message)
	if result.Error != nil {
		fmt.Printf("error retrieving user channels: %v", result.Error)
		return result.Error
	}

	// Update the state
	message.Status = status

	result = mc.db.GetNotificationsDB().WithContext(ctx).Save(&message)
	if result.Error != nil {
		fmt.Printf("error updating user channel: %v", result.Error)
		return result.Error
	}

	return nil
}

// AddMessage adds a message to the database
func (mc *MessageController) AddMessage(ctx context.Context, message *entity.Message) error {
	message.ID = uuid.New().String()
	message.Status = entity.MessagePending
	message.Version = 1

	// Insert
	result := mc.db.GetNotificationsDB().WithContext(ctx).Create(&message)
	if result.Error != nil {
		fmt.Printf("error creating user channel: %v", result.Error)
		return result.Error
	}

	return nil
}
