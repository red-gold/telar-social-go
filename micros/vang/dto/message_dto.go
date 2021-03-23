package dto

import (
	uuid "github.com/gofrs/uuid"
)

type Message struct {
	ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
	OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	RoomId      uuid.UUID `json:"roomId" bson:"roomId"`
	Text        string    `json:"text" bson:"text"`
	CreatedDate int64     `json:"createdDate" bson:"createdDate"`
	UpdatedDate int64     `json:"updatedDate" bson:"updatedDate"`
}
