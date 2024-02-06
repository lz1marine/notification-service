package controller

import (
	"encoding/json"
	"fmt"
	"time"

	apiv1 "github.com/lz1marine/notification-service/api/v1"
	channel "github.com/lz1marine/notification-service/pkg/channel"
	client "github.com/lz1marine/notification-service/pkg/client"
	"github.com/lz1marine/notification-service/pkg/queue"
)

type ChannelHandler struct {
	channel          channel.Channel
	distributedQueue queue.ReadWriter
	templateClient   client.TemplateReader
	backupClient     client.BackupMessageRemover

	maxConnections uint
}

func NewChannelHandler(ch channel.Channel, distributedQueue queue.ReadWriter, th client.TemplateReader, bmr client.BackupMessageRemover, maxConnections uint) *ChannelHandler {
	return &ChannelHandler{
		channel:          ch,
		distributedQueue: distributedQueue,
		templateClient:   th,
		backupClient:     bmr,
		maxConnections:   maxConnections,
	}
}

// Start runs a control loop that handles notifications
func (ch *ChannelHandler) Start() {
	maxConCh := make(chan struct{}, ch.maxConnections)

	for {
		maxConCh <- struct{}{}

		go func() {
			defer func() {
				<-maxConCh
			}()

			notification, err := ch.distributedQueue.Pop(ch.channel.Name())
			if err != nil {
				// TODO: log info
				// fmt.Printf("no new notifications: %v\n", err)
				time.Sleep(time.Second)
				return
			}

			var req apiv1.NotificationRequest
			err = json.Unmarshal([]byte(notification), &req)
			if err != nil {
				fmt.Printf("failed to unmarshal notification: %v\n%s\n", err, notification)
				return
			}

			err = ch.Notify(&req)
			if err != nil {
				fmt.Printf("failed to notify: %v\n", err)
				return
			}

			err = ch.backupClient.Remove(&req)
			if err != nil {
				fmt.Printf("failed to remove message: %v\n", err)
				return
			}
		}()

	}
}

// Notify handles the notification flow.
func (ch *ChannelHandler) Notify(req *apiv1.NotificationRequest) error {
	t, err := ch.templateClient.Read(req.TemplateID)
	if err != nil {
		return err
	}

	chanMessage := &queue.Message{
		Recepients: req.To,
		Message:    req.Message,
		Subject:    req.Subject,
		Template:   t,
	}

	// TODO: Log info
	fmt.Printf("message: %+v\n", chanMessage)
	return ch.channel.Notify(chanMessage)
}
