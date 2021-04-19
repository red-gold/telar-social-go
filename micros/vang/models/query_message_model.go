package models

import (
	uuid "github.com/gofrs/uuid"
)

type QueryMessageModel struct {
	ReqUserId uuid.UUID `json:"reqUserId" bson:"reqUserId"`
	RoomId    uuid.UUID `json:"roomId" bson:"roomId"`
	Page      int64     `json:"page" bson:"page"`
	Lte       int64     `json:"lte" bson:"lte"`
	Gte       int64     `json:"gte" bson:"gte"`
}
