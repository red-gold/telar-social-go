package plivo

const HOLD = "hold"
const UNHOLD = "unhold"
const HANGUP = "hangup"
const RESUME_CALL = "resume_call"
const ABORT_TRANSFER = "abort_transfer"
const VOICEMAIL_DROP = "voicemail_drop"


type MultiPartyCallMemberActionPayload struct {
	Action string `json:"action" url:"action"`
}

type MultiPartyCallMember struct {
	NodeID        string `json:"node_id" url:"node_id"`
	PhloID        string `json:"phlo_id" url:"phlo_id"`
	NodeType      string `json:"node_type" url:"node_type"`
	MemberAddress string `json:"member_address" url:"member_address"`
	BaseResource
}

func (self *MultiPartyCall) Member(memberID string) (response *MultiPartyCallMember) {
	response = &MultiPartyCallMember{self.NodeID, self.PhloID, self.NodeType, memberID, BaseResource{self.client}}
	return
}

func (self *MultiPartyCallMember) AbortTransfer() (*NodeActionResponse,error) {
	return self.update(MultiPartyCallMemberActionPayload{ABORT_TRANSFER})
}

func (service *MultiPartyCallMember) ResumeCall() (*NodeActionResponse, error) {
	return service.update(MultiPartyCallMemberActionPayload{RESUME_CALL})
}
func (service *MultiPartyCallMember) VoiceMailDrop() (*NodeActionResponse, error) {
	return service.update(MultiPartyCallMemberActionPayload{VOICEMAIL_DROP})
}
func (service *MultiPartyCallMember) HangUp() (*NodeActionResponse, error) {
	return service.update(MultiPartyCallMemberActionPayload{HANGUP})
}
func (service *MultiPartyCallMember) Hold() (*NodeActionResponse, error) {
	return service.update(MultiPartyCallMemberActionPayload{HOLD})
}
func (service *MultiPartyCallMember) UnHold() (*NodeActionResponse, error) {
	return service.update(MultiPartyCallMemberActionPayload{UNHOLD})
}

func (service *MultiPartyCallMember) update(params MultiPartyCallMemberActionPayload) (response *NodeActionResponse, err error) {
	req, err := service.client.NewRequest("POST", params, "phlo/%s/%s/%s/members/%s", service.PhloID, service.NodeType,
		service.NodeID, service.MemberAddress)
	if err != nil {
		return
	}
	response = &NodeActionResponse{}
	err = service.client.ExecuteRequest(req, response)

	return
}