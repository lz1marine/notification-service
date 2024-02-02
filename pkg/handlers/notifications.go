package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/lz1marine/notification-service/api/v1"
	"github.com/lz1marine/notification-service/pkg/adapters"
	"github.com/lz1marine/notification-service/pkg/entities"
)

// GetChannels responds with the list of all channels, including which ones are enabled
// get /v1/notifications
func GetChannels(c *gin.Context) {
	chans := entities.GetChannels()
	response := adapters.ToChannelResponse(chans)

	c.JSON(http.StatusOK, response)
}

// GetChannels returns the list of all channels, and tells the caller which ones are enabled
// get /v1/notifications/sub/:id
func GetSubNotifications(c *gin.Context) {
	userId := c.Params.ByName("id")
	chans := entities.GetUserChannel(userId)
	response := adapters.ToChannelResponse(chans)

	c.JSON(http.StatusOK, response)
}

// GetChannels returns the list of all channels, and tells the caller which ones are enabled
// patch /v1/notifications/sub/:id
func PatchSubNotifications(c *gin.Context) {
	var req apiv1.SetChannelsRequest

	userId := c.Params.ByName("id")
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chans := adapters.ToChannelEntity(&req)

	err := entities.SetUserChannel(userId, chans)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
