package plivo

type Meta struct {
	Previous *string
	Next     *string

	TotalCount int64
	Offset     int64
	Limit      int64
}

type BaseListResponse struct {
	ApiID string `json:"api_id" url:"api_id"`
	Meta  Meta   `json:"meta" url:"meta"`
}

type BaseResponse struct {
	ApiId   string `json:"api_id" url:"api_id"`
	Message string `json:"message" url:"message"`
}

func (self Application) ID() string {
	return self.AppID
}

func (self Account) ID() string {
	return self.AuthID
}

func (self Subaccount) ID() string {
	return self.AuthID
}

func (self Call) ID() string {
	return self.CallUUID
}

func (self LiveCall) ID() string {
	return self.CallUUID
}

func (self Conference) ID() string {
	return self.ConferenceName
}

func (self Endpoint) ID() string {
	return self.EndpointID
}

func (self Message) ID() string {
	return self.MessageUUID
}

func (self Number) ID() string {
	return self.Number
}

func (self PhoneNumber) ID() string {
	return self.Number
}

func (self Pricing) ID() string {
	return self.CountryISO
}

func (self Recording) ID() string {
	return self.RecordingID
}
