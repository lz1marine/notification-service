package handlers

import (
	"fmt"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/lz1marine/notification-service/api/v1"
	"github.com/lz1marine/notification-service/pkg/adapters"
	"github.com/lz1marine/notification-service/pkg/clients"
	"github.com/lz1marine/notification-service/pkg/entities"
	"github.com/lz1marine/notification-service/pkg/queue"

	"github.com/google/uuid"
)

type NotificationHandler struct {
	distQueue queue.Writer
	backup    clients.BackupMessageSender
}

func NewNotificationHandler(distQueue queue.Writer, backup clients.BackupMessageSender) *NotificationHandler {
	return &NotificationHandler{
		distQueue: distQueue,
		backup:    backup,
	}
}

// PostNotification posts a notification to a channel
// post /v1/internal/notifications
func (nh *NotificationHandler) PostNotification(c *gin.Context) {
	eventId := c.Params.ByName("id")

	var req apiv1.ChannelNotificationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := adapters.ToMessageEntity(&req)

	// TODO: the following two should be a transaction
	err := entities.AddMessage(message, eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = entities.AddMessageTopic(message, req.Topic)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emails := entities.GetEmails(req.Topic)

	// Save to the backup table
	queueReq := &apiv1.NotificationRequest{
		Subject:    req.Subject,
		Message:    req.Message,
		TemplateID: req.TemplateID,
	}

	// Push to workers
	err = nh.pushToWorker(queueReq, req.Channel, emails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// pushToWorker Pushes all notifications to the workers
// TODO: maybe fanout here
func (nh *NotificationHandler) pushToWorker(req *apiv1.NotificationRequest, channel string, emails []string) error {
	for _, email := range emails {
		if !isValidEmail(email) {
			continue
		}

		req.To = []string{email}
		req.ID = uuid.New().String()

		// Push to backup table for reevaluation
		err := nh.backup.Send(req)
		if err != nil {
			return err
		}

		// Push to workers
		err = nh.distQueue.Push(req, channel)
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
