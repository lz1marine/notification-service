package channels

import apiv1 "github.com/lz1marine/notification-service/api/v1"

type Channel interface {
	Name() string
	Notify(req apiv1.NotificationRequest) error
}
