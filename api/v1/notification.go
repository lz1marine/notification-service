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
