package queue

import apiv1 "github.com/lz1marine/notification-service/api/v1"

// Reader is used to read messages from a queue
type Reader interface {
	// Pop returns and removes the next message from the queue
	Pop(channel string) (string, error)
}

// Writer is used to write messages to a queue
type Writer interface {
	// Push adds a message to the queue
	Push(req *apiv1.NotificationRequest, channel string) error
}

// ReadWriter is used to read and write messages to a queue
type ReadWriter interface {
	Reader
	Writer
}
