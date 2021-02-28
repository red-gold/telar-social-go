package models

import (
	uuid "github.com/gofrs/uuid"
)

type CommentModel struct {
	ObjectId         uuid.UUID `json:"objectId"`
	Score            int64     `json:"score"`
	OwnerUserId      uuid.UUID `json:"ownerUserId"`
	OwnerDisplayName string    `json:"ownerDisplayName"`
	OwnerAvatar      string    `json:"ownerAvatar"`
	PostId           uuid.UUID `json:"postId"`
	Text             string    `json:"text"`
	Deleted          bool      `json:"deleted"`
	DeletedDate      int64     `json:"deletedDate"`
	CreatedDate      int64     `json:"created_date"`
	LastUpdated      int64     `json:"last_updated"`
}
