package dto

import (
	uuid "github.com/gofrs/uuid"
)

type Vote struct {
	ObjectId         uuid.UUID `json:"objectId" bson:"objectId"`
	OwnerUserId      uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	OwnerDisplayName string    `json:"ownerDisplayName" bson:"ownerDisplayName"`
	OwnerAvatar      string    `json:"ownerAvatar" bson:"ownerAvatar"`
	PostId           uuid.UUID `json:"postId" bson:"postId"`
	TypeId           int       `json:"type" bson:"type"`
	CreatedDate      int64     `json:"created_date" bson:"created_date"`
}
