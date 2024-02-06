package channel

import "github.com/lz1marine/notification-service/pkg/queue"

// Channel is a channel for messages
type Channel interface {
	// Name returns the name of the channel
	Name() string

	// Notify sends a notification
	Notify(req *queue.Message) error
}
