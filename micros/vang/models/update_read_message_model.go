package models

import (
	uuid "github.com/gofrs/uuid"
)

type UpdateReadMessageModel struct {
	RoomId             uuid.UUID `json:"roomId" bson:"roomId"`
	Amount             int64     `json:"amount" bson:"amount"`
	MessageCreatedDate int64     `json:"messageCreatedDate" bson:"messageCreatedDate"`
	MessageId          uuid.UUID `json:"messageId" bson:"messageId"`
}
