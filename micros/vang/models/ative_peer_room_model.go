package models

import "github.com/gofrs/uuid"

type ActivePeerRoomModel struct {
	PeerUserId uuid.UUID `json:"peerUser" bson:"peerUser"`
}
