package models

import uuid "github.com/gofrs/uuid"

type RoomMemberModel struct {
	ObjectId uuid.UUID `json:"objectId"`
	FullName string    `json:"fullName"`
	Avatar   string    `json:"avatar"`
	LastSeen int64     `json:"lastSeen"`
}
