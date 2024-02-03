package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"

	apiv1 "github.com/lz1marine/notification-service/api/v1"
	"github.com/lz1marine/notification-service/pkg/channels"
	"github.com/lz1marine/notification-service/pkg/entities"
	"github.com/lz1marine/notification-service/pkg/queue"
)

type ChannelHandler struct {
	channel   channels.Channel
	distQueue queue.ReadWriter
}

func NewChannelHandler(ch channels.Channel, distQueue queue.ReadWriter) *ChannelHandler {
	return &ChannelHandler{
		channel:   ch,
		distQueue: distQueue,
	}
}

// Start runs a control loop that handles notifications
func (ch *ChannelHandler) Start() {
	for {
		notification, err := ch.distQueue.Pop(ch.channel.Name())
		if err != nil {
			fmt.Printf("failed to pop notification: %v\n", err)
			continue
		}

		var req apiv1.NotificationRequest
		err = json.Unmarshal([]byte(notification), &req)
		if err != nil {
			fmt.Printf("failed to unmarshal notification: %v\n%s\n", err, notification)
			continue
		}

		ch.Notify(req)
	}
}

// Notify handles the notification flow. First, it tries to
func (ch *ChannelHandler) Notify(req apiv1.NotificationRequest) {
	message := entities.GetMessage(req.MessageID)
	entities.SetMessageStatus(req.MessageID, entities.Active)

	topic := entities.GetTopic(req.MessageID)
	emails := entities.GetEmails(topic)

	// TODO: Log info
	fmt.Printf("message: %+v\n", message)

	chanMessage := &queue.Message{
		Recepients: emails,
		Message:    message.Message,
		Title:      message.Title,
		Template:   generateTemplate(message.TemplateID),
	}

	err := ch.channel.Notify(chanMessage)
	if err != nil {
		entities.SetMessageStatus(req.MessageID, entities.Failed)
	}

	entities.SetMessageStatus(req.MessageID, entities.Sent)
}

// TODO: not yet working, maybe move to another logical place
func generateTemplate(templateID *string) *template.Template {
	if templateID == nil {
		return nil
	}

	tmp := entities.GetTemplates(*templateID)
	res, err := template.New("template").Parse(tmp.Template)
	if err != nil {
		fmt.Printf("failed to parse template: %v", err)
	}

	return res
}
