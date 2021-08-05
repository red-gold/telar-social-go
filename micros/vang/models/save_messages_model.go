package models

import (
	uuid "github.com/gofrs/uuid"
)

type SaveMessagesModel struct {
	UserId         uuid.UUID      `json:"userId" bson:"userId"`
	RoomId         uuid.UUID      `json:"roomId" bson:"roomId"`
	Messages       []MessageModel `json:"messages" bson:"messages"`
	DeactivePeerId uuid.UUID      `json:"deactivePeerId" bson:"deactivePeerId"`
}
