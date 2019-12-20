package utils

import (
	"fmt"

	"github.com/plivo/plivo-go"
)

type PhoneClient struct {
	client       *plivo.Client
	sourceNumber string
}

// NewPhone initialize phone config
func NewPhone(authId string, authToken string, sourceNumber string) (*PhoneClient, error) {
	fmt.Println("Start Phone client initializing...")
	client, err := plivo.NewClient(authId, authToken, &plivo.ClientOptions{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Phone client initialized.")
	return &PhoneClient{client: client, sourceNumber: sourceNumber}, nil
}

// SendSms send sms
func (pc *PhoneClient) SendSms(phoneNumber string, message string) ([]string, error) {
	fmt.Println("Start sending message...")

	req := plivo.MessageCreateParams{
		Src:  pc.sourceNumber,
		Dst:  phoneNumber,
		Text: message,
	}
	res, err := pc.client.Messages.Create(req)
	return res.MessageUUID, err
}
