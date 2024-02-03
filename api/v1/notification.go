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
	MessageID string `json:"message_id"`
}

type ChannelNotificationRequest struct {
	Channel    string  `json:"channel"`
	Message    string  `json:"message"`
	Topic      string  `json:"topic_id"`
	Title      *string `json:"title,omitempty"`
	TemplateID *string `json:"template_id,omitempty"`
}
