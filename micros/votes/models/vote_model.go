package models

import (
	uuid "github.com/gofrs/uuid"
)

type VoteModel struct {
	ObjectId         uuid.UUID `json:"objectId"`
	OwnerUserId      uuid.UUID `json:"ownerUserId"`
	OwnerDisplayName string    `json:"ownerDisplayName"`
	OwnerAvatar      string    `json:"ownerAvatar"`
	PostId           uuid.UUID `json:"postId"`
	TypeId           int       `json:"type"`
	CreatedDate      int64     `json:"created_date"`
}
