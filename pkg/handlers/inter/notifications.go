package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/lz1marine/notification-service/api/v1"
	"github.com/lz1marine/notification-service/pkg/adapters"
	"github.com/lz1marine/notification-service/pkg/entities"
	"github.com/lz1marine/notification-service/pkg/queue"
)

type NotificationHandler struct {
	distQueue queue.Writer
}

func NewNotificationHandler(distQueue queue.Writer) *NotificationHandler {
	return &NotificationHandler{
		distQueue: distQueue,
	}
}

// PostNotification posts a notification to our channel
// post /v1/internal/notifications
func (nh *NotificationHandler) PostNotification(c *gin.Context) {
	var req apiv1.ChannelNotificationRequest

	eventId := c.Params.ByName("id")
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

	// TODO: push to queue
	nh.distQueue.Push(message.ID, req.Channel)

	c.JSON(http.StatusOK, nil)
}
