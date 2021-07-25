package models

import "github.com/gofrs/uuid"

type ActivePeerRoomModel struct {
	PeerUserId         uuid.UUID `json:"peerUserId" bson:"peerUserId"`
	SocialName         string    `json:"socialName" bson:"socialName"`
	ResponseActionType string    `json:"responseActionType" bson:"responseActionType"`
}
