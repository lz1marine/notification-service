package channel

import (
	"fmt"

	"github.com/lz1marine/notification-service/pkg/queue"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type smsChannel struct {
	sender string
	client *twilio.RestClient
}

// NewSMSChannel creates a new instance of SMSChannel
func NewSMSChannel(acc, token, sender string) *smsChannel {
	params := twilio.ClientParams{
		Username: acc,
		Password: token,
	}

	client := twilio.NewRestClientWithParams(params)
	return &smsChannel{
		sender: sender,
		client: client,
	}
}

// Name returns the name of the channel
func (sms *smsChannel) Name() string {
	return "sms"
}

// Notify sends a notification
func (sms *smsChannel) Notify(m *queue.Message) error {
	for _, recipient := range m.Recepients {
		params := &openapi.CreateMessageParams{}
		params.SetTo(recipient)

		fmt.Printf("sms.sender: %s\n", sms.sender)
		params.SetFrom(sms.sender)
		params.SetBody(m.Message)

		_, err := sms.client.Api.CreateMessage(params)
		if err != nil {
			fmt.Println(err.Error())
		}

	}

	return nil
}
