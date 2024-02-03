package adapters

import (
	apiv1 "github.com/lz1marine/notification-service/api/v1"
	v1 "github.com/lz1marine/notification-service/api/v1"
	"github.com/lz1marine/notification-service/pkg/entities"
)

func ToChannelResponse(channel []entities.Channels) *apiv1.ChannelResponse {
	chans := make([]v1.Channel, 0, len(channel))
	for _, v := range channel {
		chans = append(chans, apiv1.Channel{
			Name:      v.Name,
			IsEnabled: v.IsEnabled,
		})
	}

	return &v1.ChannelResponse{
		ChannelWrapper: apiv1.ChannelWrapper{
			Channels: chans,
		},
	}
}

func ToChannelEntity(channel *apiv1.SetChannelsRequest) []entities.Channels {
	chans := make([]entities.Channels, 0, len(channel.ChannelWrapper.Channels))
	for _, v := range channel.ChannelWrapper.Channels {
		chans = append(chans, entities.Channels{
			Name:      v.Name,
			IsEnabled: v.IsEnabled,
		})
	}

	return chans
}

func ToNotificationEntity(notification *apiv1.ChannelNotificationRequest) *entities.Notifications {
	return &entities.Notifications{
		Title:      *notification.Title,
		Message:    notification.Message,
		TemplateID: *notification.TemplateID,
		ChannelID:  entities.GetChannelID(notification.Channel).ID,
	}
}

// TODO: remove
func ToTempNotificationEntity(notification *apiv1.NotificationRequest) *entities.Notifications {
	res := &entities.Notifications{
		Message: notification.Message,
	}

	if notification.Title != nil {
		res.Title = *notification.Title
	}

	if notification.TemplateID != nil {
		res.TemplateID = *notification.TemplateID
	}

	return res
}
