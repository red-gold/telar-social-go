package plivo

type MultiPartyCall struct {
	Node
	BaseResource
}

type MultiPartyCallActionPayload struct {
	Action        string `json:"action" url:"action"`
	To            string `json:"to" url:"to"`
	Role          string `json:"role" url:"role"`
	TriggerSource string `json:"trigger_source" url:"trigger_source"`
}


func (self *MultiPartyCall) update(params MultiPartyCallActionPayload) (response *NodeActionResponse, err error) {
	req, err := self.client.NewRequest("POST", params, "phlo/%s/%s/%s", self.PhloID, self.NodeType,
		self.NodeID)
	if (err != nil) {
		return
	}
	response = &NodeActionResponse{}
	err = self.client.ExecuteRequest(req, response)

	return
}

func (self *MultiPartyCall) Call(params MultiPartyCallActionPayload) (*NodeActionResponse, error) {
	return self.update(params)
}

func (self *MultiPartyCall) WarmTransfer(params MultiPartyCallActionPayload) (response *NodeActionResponse,
	err error) {
	return self.update(params)
}

func (self *MultiPartyCall) ColdTransfer(params MultiPartyCallActionPayload) (response *NodeActionResponse,
	err error) {
	return self.update(params)
}
