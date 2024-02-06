package v1

type SetChannelsRequest struct {
	ChannelWrapper
}

type SetChannelsResponse struct {
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
	ID         string   `json:"id"`
	To         []string `json:"to"`
	Subject    *string  `json:"subject,omitempty"`
	TemplateID *string  `json:"templateId,omitempty"`
	Message    string   `json:"message"`
}

type ChannelNotificationRequest struct {
	Channel    string  `json:"channel"`
	Message    string  `json:"message"`
	Topic      string  `json:"topic_id"`
	Subject    *string `json:"subject,omitempty"`
	TemplateID *string `json:"template_id,omitempty"`
}

type ChannelNotificationResponse struct {
}
