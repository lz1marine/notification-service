package v1

type SetChannelsRequest struct {
	ChannelWrapper
}

type ChannelResponse struct {
	ChannelWrapper
}

type ChannelWrapper struct {
	Channels []Channel `json:"channels"`
}

type Channel struct {
	Name      string `json:"name"`
	IsEnabled bool   `json:"is_enabled"`
}

type NotificationRequest struct {
	Message    string  `json:"message"`
	Title      *string `json:"title,omitempty"`
	TemplateID *string `json:"template_id,omitempty"`

	Topic string `json:"topic_id"`
}

type ChannelNotificationRequest struct {
	NotificationRequest

	Channel string `json:"channel"`
}
