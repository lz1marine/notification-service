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
	"github.com/lz1marine/notification-service/pkg/database/entity"
	"github.com/lz1marine/notification-service/pkg/queue"
	"github.com/swaggo/swag/example/celler/httputil"

	"github.com/google/uuid"
)

// InternalHandler is the handler for internal user notifications
type InternalHandler struct {
	distQueue         queue.Writer
	backup            client.BackupMessageSender
	channelController *controller.ChannelController
	messageController *controller.MessageController
	userController    *controller.UserController
}

// NewNotificationHandler creates a new NotificationHandler
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
// @Summary      Post a notification to a channel
// @Description  post a notification to a channel given an event id
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        id path string true "Event ID"
// @Param 		 request body v1.ChannelNotificationRequest true "The request body."
// @Success      200  {object}  v1.ChannelNotificationResponse
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/v1/internal/notifications/{id} [post]
func (ih *InternalHandler) PostNotification(c *gin.Context) {
	eventId := c.Params.ByName("id")

	var req apiv1.ChannelNotificationRequest
	if err := c.BindJSON(&req); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	chanId, err := ih.channelController.GetChannelID(context.Background(), req.Channel)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	message := adapter.ToMessageEntity(&req, chanId, eventId)

	err = ih.messageController.AddMessage(context.Background(), message)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	profiles, err := ih.userController.GetProfiles(context.Background(), req.Topic)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	// Save to the backup table
	queueReq := &apiv1.NotificationRequest{
		Subject:    req.Subject,
		Message:    req.Message,
		TemplateID: req.TemplateID,
	}

	// Push to workers
	err = ih.pushToWorker(queueReq, req.Channel, profiles)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, apiv1.ChannelNotificationResponse{})
}

// pushToWorker Pushes all notifications to the workers
// TODO: maybe fanout here
func (ih *InternalHandler) pushToWorker(req *apiv1.NotificationRequest, channel string, recepients []entity.Profile) error {
	for _, r := range recepients {
		recepient, ok := getValidRecepient(channel, r)
		if !ok {
			continue
		}

		req.To = []string{recepient}
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
func getValidRecepient(channel string, recepient entity.Profile) (string, bool) {
	res := ""
	switch channel {
	case "email":
		res = recepient.Email
		_, err := mail.ParseAddress(res)
		if err != nil {
			fmt.Printf("failed to parse email %s: %v", res, err)
			return "", false
		}
	case "sms":
		// TODO: add phone number validation
		res = recepient.Phone
	}

	return res, true
}
