package plivo

type MessageService struct {
	client *Client
}

type MessageCreateParams struct {
	Src  string `json:"src,omitempty" url:"src,omitempty"`
	Dst  string `json:"dst,omitempty" url:"dst,omitempty"`
	Text string `json:"text,omitempty" url:"text,omitempty"`
	// Optional parameters.
	Type      string      `json:"type,omitempty" url:"type,omitempty"`
	URL       string      `json:"url,omitempty" url:"url,omitempty"`
	Method    string      `json:"method,omitempty" url:"method,omitempty"`
	Trackable bool        `json:"trackable,omitempty" url:"trackable,omitempty"`
	Log       interface{} `json:"log,omitempty" url:"log,omitempty"`
	// Either one of src and powerpackuuid should be given
	PowerpackUUID string `json:"powerpack_uuid,omitempty" url:"powerpack_uuid,omitempty"`
}

type Message struct {
	ToNumber         string `json:"to_number,omitempty" url:"to_number,omitempty"`
	FromNumber       string `json:"from_number,omitempty" url:"from_number,omitempty"`
	CloudRate        string `json:"cloud_rate,omitempty" url:"cloud_rate,omitempty"`
	MessageType      string `json:"message_type,omitempty" url:"message_type,omitempty"`
	ResourceURI      string `json:"resource_uri,omitempty" url:"resource_uri,omitempty"`
	CarrierRate      string `json:"carrier_rate,omitempty" url:"carrier_rate,omitempty"`
	MessageDirection string `json:"message_direction,omitempty" url:"message_direction,omitempty"`
	MessageState     string `json:"message_state,omitempty" url:"message_state,omitempty"`
	TotalAmount      string `json:"total_amount,omitempty" url:"total_amount,omitempty"`
	MessageUUID      string `json:"message_uuid,omitempty" url:"message_uuid,omitempty"`
	MessageTime      string `json:"message_time,omitempty" url:"message_time,omitempty"`
}

// Stores response for ending a message.
type MessageCreateResponseBody struct {
	Message     string   `json:"message" url:"message"`
	ApiID       string   `json:"api_id" url:"api_id"`
	MessageUUID []string `json:"message_uuid" url:"message_uuid"`
	Error       string   `json:"error" url:"error"`
}

type MessageList struct {
	BaseListResponse
	Objects []Message `json:"objects" url:"objects"`
}

type MessageListParams struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

func (service *MessageService) List(params MessageListParams) (response *MessageList, err error) {
	req, err := service.client.NewRequest("GET", params, "Message")
	if err != nil {
		return
	}
	response = &MessageList{}
	err = service.client.ExecuteRequest(req, response)
	return
}

func (service *MessageService) Get(messageUuid string) (response *Message, err error) {
	req, err := service.client.NewRequest("GET", nil, "Message/%s", messageUuid)
	if err != nil {
		return
	}
	response = &Message{}
	err = service.client.ExecuteRequest(req, response)
	return
}

func (service *MessageService) Create(params MessageCreateParams) (response *MessageCreateResponseBody, err error) {
	req, err := service.client.NewRequest("POST", params, "Message")
	if err != nil {
		return
	}
	response = &MessageCreateResponseBody{}
	err = service.client.ExecuteRequest(req, response)
	return
}
