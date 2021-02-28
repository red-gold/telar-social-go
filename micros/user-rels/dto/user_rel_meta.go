package dto

import uuid "github.com/satori/go.uuid"

type UserRelMeta struct {
	UserId   uuid.UUID `json:"userId" bson:"userId"`
	FullName string    `json:"fullName" bson:"fullName"`
	Avatar   string    `json:"avatar" bson:"avatar"`
}
