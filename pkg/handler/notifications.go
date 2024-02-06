package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/lz1marine/notification-service/api/v1"
	"github.com/lz1marine/notification-service/pkg/adapter"
	"github.com/lz1marine/notification-service/pkg/database/controller"
	"github.com/lz1marine/notification-service/pkg/database/entity"

	"github.com/swaggo/swag/example/celler/httputil"
)

// ExternalHandler is the handler for external user notifications
type ExternalHandler struct {
	channelController *controller.ChannelController
}

// NewExternalHandler creates a new ExternalHandler
func NewExternalHandler(channelController *controller.ChannelController) *ExternalHandler {
	return &ExternalHandler{
		channelController: channelController,
	}
}

// GetChannels responds with the list of all channels, including whether they are enabled
// PostNotification posts a notification to a channel
// @Summary      Gets the list of all channels, including whether they are enabled
// @Description  gets the list of all channels, including whether they are enabled
// @Tags         notifications
// @Produce      json
// @Success      200  {object}  v1.ChannelResponse
// @Router       /api/v1/notifications [get]
func (eh *ExternalHandler) GetChannels(c *gin.Context) {
	chans := eh.channelController.GetAllChannels(context.Background())
	response := adapter.ToChannelResponse(chans)

	c.JSON(http.StatusOK, response)
}

// GetSubNotifications returns the list of all channels subscribed to by the user, including whether they are enabled for the user
// @Summary      Gets the list of all channels subscribed to by the user, including whether they are enabled for the user
// @Description  gets the list of all channels subscribed to by the user, including whether they are enabled for the user
// @Tags         notifications
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200  {object}  v1.ChannelResponse
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/v1/notifications/sub/{id} [get]
func (eh *ExternalHandler) GetSubNotifications(c *gin.Context) {
	userId := c.Params.ByName("id")
	chans, err := eh.channelController.GetUserChannels(context.Background(), userId)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
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

// PatchSubNotifications subscribes and unsubscribes a user to the given channels
// @Summary      Patch the channel list that the user has subscribed to
// @Description  patch he channel list that the user has subscribed to
// @Tags         notifications
// @Produce      json
// @Param        id path string true "User ID"
// @Param 		 request body v1.SetChannelsRequest true "The request body."
// @Success      200  {object}  v1.SetChannelsResponse
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/v1/notifications/sub/{id} [patch]
func (eh *ExternalHandler) PatchSubNotifications(c *gin.Context) {
	var req apiv1.SetChannelsRequest

	userId := c.Params.ByName("id")
	if err := c.BindJSON(&req); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	chans := adapter.ToChannelEntity(&req)

	err := eh.channelController.SetUserChannels(context.Background(), userId, chans)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
