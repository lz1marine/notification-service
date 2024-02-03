package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/lz1marine/notification-service/api/v1"
	"github.com/lz1marine/notification-service/pkg/adapters"
	"github.com/lz1marine/notification-service/pkg/entities"
)

// PostNotification posts a notification to our channel
// post /v1/internal/notifications
func PostNotification(c *gin.Context) {
	var req apiv1.ChannelNotificationRequest

	eventId := c.Params.ByName("id")
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := adapters.ToNotificationEntity(&req)

	err := entities.AddNotification(notification, eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: push to queue

	c.JSON(http.StatusOK, nil)
}
