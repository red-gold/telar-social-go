package dto

import uuid "github.com/gofrs/uuid"

type UserRelMeta struct {
	UserId   uuid.UUID `json:"userId" bson:"userId"`
	FullName string    `json:"fullName" bson:"fullName"`
	Avatar   string    `json:"avatar" bson:"avatar"`
}
