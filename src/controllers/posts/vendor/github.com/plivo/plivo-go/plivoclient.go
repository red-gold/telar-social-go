package plivo

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"runtime"
	)

const baseUrlString = "https://api.plivo.com/"
const baseRequestString = "/v1/Account/%s/"


type Client struct {
	BaseClient

	Messages     *MessageService
	Accounts     *AccountService
	Subaccounts  *SubaccountService
	Applications *ApplicationService
	Endpoints    *EndpointService
	Numbers      *NumberService
	PhoneNumbers *PhoneNumberService
	Pricing      *PricingService // TODO Rename?
	Recordings   *RecordingService
	Calls        *CallService
	LiveCalls    *LiveCallService
	QueuedCalls  *QueuedCallService
	Conferences  *ConferenceService
}

/*
To set a proxy for all requests, configure the Transport for the HttpClient passed in:

	&http.Client{
 		Transport: &http.Transport{
 			Proxy: http.ProxyURL("http//your.proxy.here"),
 		},
 	}

Similarly, to configure the timeout, set it on the HttpClient passed in:

	&http.Client{
 		Timeout: time.Minute,
 	}
*/
func NewClient(authId, authToken string, options *ClientOptions) (client *Client, err error) {

	client = &Client{}

	if len(authId) == 0 {
		authId = os.Getenv("PLIVO_AUTH_ID")
	}
	if len(authToken) == 0 {
		authToken = os.Getenv("PLIVO_AUTH_TOKEN")
	}
	client.AuthId = authId
	client.AuthToken = authToken
	client.userAgent = fmt.Sprintf("%s/%s (Go: %s)", "plivo-go", sdkVersion, runtime.Version())

	baseUrl, err := url.Parse(baseUrlString) // Todo: handle error case?

	client.BaseUrl = baseUrl
	client.httpClient = &http.Client{
		Timeout: time.Minute,
	}

	if options.HttpClient != nil {
		client.httpClient = options.HttpClient
	}

	client.Messages = &MessageService{client: client}
	client.Accounts = &AccountService{client: client}
	client.Subaccounts = &SubaccountService{client: client}
	client.Applications = &ApplicationService{client: client}
	client.Endpoints = &EndpointService{client: client}
	client.Numbers = &NumberService{client: client}
	client.PhoneNumbers = &PhoneNumberService{client: client}
	client.Pricing = &PricingService{client: client}
	client.Recordings = &RecordingService{client: client}
	client.Calls = &CallService{client: client}
	client.LiveCalls = &LiveCallService{client: client}
	client.QueuedCalls = &QueuedCallService{client: client}
	client.Conferences = &ConferenceService{client: client}

	return
}

func (client *Client) NewRequest(method string, params interface{}, formatString string,
formatParams ...interface{}) (*http.Request, error) {
	formatParams = append([]interface{}{client.AuthId}, formatParams...)
	formatString = fmt.Sprintf("%s/%s", "%s", formatString)
	return client.BaseClient.NewRequest(method, params, baseRequestString, formatString, formatParams...)
}
