package channel

import "github.com/lz1marine/notification-service/pkg/queue"

type Channel interface {
	Name() string
	Notify(req *queue.Message) error
}
