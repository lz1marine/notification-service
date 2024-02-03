package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/lz1marine/notification-service/api/v1"
	"github.com/lz1marine/notification-service/pkg/adapters"
	"github.com/lz1marine/notification-service/pkg/channels"
	"github.com/lz1marine/notification-service/pkg/entities"
)

type ChannelHandler struct {
	channel channels.Channel
}

func NewChannelHandler(ch channels.Channel) *ChannelHandler {
	return &ChannelHandler{
		channel: ch,
	}
}

// PostNotification posts a notification to our channel
// post /v1/temp/notifications/:id
func (ch *ChannelHandler) PostNotification(c *gin.Context) {
	var req apiv1.NotificationRequest

	messageID := c.Params.ByName("id")
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := adapters.ToTempNotificationEntity(&req)

	err := entities.PatchNotification(notification, messageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("notification: %+v\n", notification)

	err = ch.channel.Notify(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
