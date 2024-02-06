package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/lz1marine/notification-service/api/v1"
	"github.com/lz1marine/notification-service/pkg/adapter"
	"github.com/lz1marine/notification-service/pkg/client"
	"github.com/lz1marine/notification-service/pkg/database/controller"
	"github.com/lz1marine/notification-service/pkg/queue"

	"github.com/google/uuid"
)

type InternalHandler struct {
	distQueue         queue.Writer
	backup            client.BackupMessageSender
	channelController *controller.ChannelController
	messageController *controller.MessageController
	userController    *controller.UserController
}

func NewNotificationHandler(
	distQueue queue.Writer,
	backup client.BackupMessageSender,
	channelController *controller.ChannelController,
	messageController *controller.MessageController,
	userController *controller.UserController) *InternalHandler {
	return &InternalHandler{
		distQueue:         distQueue,
		backup:            backup,
		channelController: channelController,
		messageController: messageController,
		userController:    userController,
	}
}

// PostNotification posts a notification to a channel
// post /v1/internal/notifications
func (ih *InternalHandler) PostNotification(c *gin.Context) {
	eventId := c.Params.ByName("id")

	var req apiv1.ChannelNotificationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chanId, err := ih.channelController.GetChannelID(context.Background(), req.Channel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := adapter.ToMessageEntity(&req, chanId, eventId)

	err = ih.messageController.AddMessage(context.Background(), message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emails, err := ih.userController.GetEmails(context.Background(), req.Topic)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save to the backup table
	queueReq := &apiv1.NotificationRequest{
		Subject:    req.Subject,
		Message:    req.Message,
		TemplateID: req.TemplateID,
	}

	// Push to workers
	err = ih.pushToWorker(queueReq, req.Channel, emails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// pushToWorker Pushes all notifications to the workers
// TODO: maybe fanout here
func (ih *InternalHandler) pushToWorker(req *apiv1.NotificationRequest, channel string, emails []string) error {
	for _, email := range emails {
		if !isValidEmail(email) {
			continue
		}

		req.To = []string{email}
		req.ID = uuid.New().String()

		// Push to backup table for reevaluation
		err := ih.backup.Send(req)
		if err != nil {
			return err
		}

		// Push to workers
		err = ih.distQueue.Push(req, channel)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: move to helpers
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		fmt.Printf("failed to parse email %s: %v", email, err)
		return false
	}

	return true
}
