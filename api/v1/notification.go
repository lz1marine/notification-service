package v1

// SetChannelsRequest is the request to set channels
type SetChannelsRequest struct {
	ChannelWrapper
}

// SetChannelsResponse is the response after setting channels
type SetChannelsResponse struct {
}

// ChannelResponse is the response to a channel
type ChannelResponse struct {
	ChannelWrapper
}

// ChannelWrapper is an array of channel objects
type ChannelWrapper struct {
	// Channels is an array of channel objects
	Channels []Channel `json:"channels"`
}

// Channel is a channel
type Channel struct {
	// Name is the name of the channel
	Name string `json:"name"`

	// IsEnabled is a flag indicating if the channel is active
	IsEnabled bool `json:"is_enabled"`
}

// NotificationRequest is the request to a notification
type NotificationRequest struct {
	// ID is the id of the notification
	ID string `json:"id"`

	// To is an array of recipients
	To []string `json:"to"`

	// Subject is the subject of the message
	Subject *string `json:"subject,omitempty"`

	// TemplateID is the id of the template
	TemplateID *string `json:"templateId,omitempty"`

	// Message is the content of the message
	Message string `json:"message"`
}

// ChannelNotificationRequest is the request to a channel notification
type ChannelNotificationRequest struct {
	// Channel is the name of the channel
	Channel string `json:"channel"`

	// Message is the content of the message
	Message string `json:"message"`

	// Topic is the id of the topic
	Topic string `json:"topic_id"`

	// Subject is the subject of the message
	Subject *string `json:"subject,omitempty"`

	// TemplateID is the id of the template
	TemplateID *string `json:"template_id,omitempty"`
}

// ChannelNotificationResponse is the response to a channel notification
type ChannelNotificationResponse struct {
}
