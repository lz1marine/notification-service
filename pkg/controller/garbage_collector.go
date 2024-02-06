package controller

import (
	"context"
	"time"

	client "github.com/lz1marine/notification-service/pkg/client"
	"github.com/lz1marine/notification-service/pkg/queue"
)

// GarbageCollector is the controller for garbage collection, which hanldes message state transfers, as well as message retries.
type GarbageCollector struct {
	backupClient     client.BackupMessageReadWriter
	distributedQueue queue.ReadWriter
}

// NewGarbageCollector creates a new GarbageCollector
func NewGarbageCollector(bmr client.BackupMessageReadWriter, distributedQueue queue.ReadWriter) *GarbageCollector {
	return &GarbageCollector{
		backupClient:     bmr,
		distributedQueue: distributedQueue,
	}
}

// Start runs a control loop that handles garbage collection
func (gc *GarbageCollector) Start() {
	for {
		time.Sleep(10 * time.Second)

		keys, err := gc.backupClient.GetIdleKeys(context.Background(), 5*time.Minute)
		if err != nil {
			continue
		}

		// Retry keys
		for _, key := range keys {
			resend, err := gc.backupClient.GetVal(key)
			if err != nil {
				continue
			}

			// Resend to queue
			err = gc.distributedQueue.Push(resend, resend.Channel)
			if err != nil {
				continue
			}

			// Reset the timeout
			err = gc.backupClient.Send(resend)
			if err != nil {
				continue
			}
		}
	}
}
