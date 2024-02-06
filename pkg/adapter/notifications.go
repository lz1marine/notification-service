package adapter

import (
	apiv1 "github.com/lz1marine/notification-service/api/v1"
	v1 "github.com/lz1marine/notification-service/api/v1"
	"github.com/lz1marine/notification-service/pkg/database/entity"
)

func ToChannelResponse(channel []entity.Channel) *apiv1.ChannelResponse {
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

func ToChannelEntity(channel *apiv1.SetChannelsRequest) []entity.Channel {
	chans := make([]entity.Channel, 0, len(channel.ChannelWrapper.Channels))
	for _, v := range channel.ChannelWrapper.Channels {
		chans = append(chans, entity.Channel{
			Name:      v.Name,
			IsEnabled: v.IsEnabled,
		})
	}

	return chans
}

func ToMessageEntity(notification *apiv1.ChannelNotificationRequest, channelID, eventID string) *entity.Message {
	return &entity.Message{
		Subject:    notification.Subject,
		Message:    notification.Message,
		TemplateID: notification.TemplateID,
		ChannelID:  channelID,
		EventID:    eventID,
	}
}
