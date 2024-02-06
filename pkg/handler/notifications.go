package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/lz1marine/notification-service/api/v1"
	"github.com/lz1marine/notification-service/pkg/adapter"
	"github.com/lz1marine/notification-service/pkg/database/controller"
	"github.com/lz1marine/notification-service/pkg/database/entity"
)

type ExternalHandler struct {
	channelController *controller.ChannelController
}

func NewExternalHandler(channelController *controller.ChannelController) *ExternalHandler {
	return &ExternalHandler{
		channelController: channelController,
	}
}

// GetChannels responds with the list of all channels, including which ones are enabled
// get /v1/notifications
func (eh *ExternalHandler) GetChannels(c *gin.Context) {
	chans := eh.channelController.GetAllChannels(context.Background())
	response := adapter.ToChannelResponse(chans)

	c.JSON(http.StatusOK, response)
}

// GetChannels returns the list of all channels, and tells the caller which ones are enabled
// get /v1/notifications/sub/:id
func (eh *ExternalHandler) GetSubNotifications(c *gin.Context) {
	userId := c.Params.ByName("id")
	chans, err := eh.channelController.GetUserChannels(context.Background(), userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract channels
	channels := make([]entity.Channel, 0, len(chans))
	for _, uc := range chans {
		channels = append(channels, *uc.Channel)
	}

	response := adapter.ToChannelResponse(channels)
	c.JSON(http.StatusOK, response)
}

// GetChannels returns the list of all channels, and tells the caller which ones are enabled
// patch /v1/notifications/sub/:id
func (eh *ExternalHandler) PatchSubNotifications(c *gin.Context) {
	var req apiv1.SetChannelsRequest

	userId := c.Params.ByName("id")
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chans := adapter.ToChannelEntity(&req)

	err := eh.channelController.SetUserChannels(context.Background(), userId, chans)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
