package dto

import (
	uuid "github.com/gofrs/uuid"
)

type Comment struct {
	ObjectId         uuid.UUID `json:"objectId" bson:"objectId"`
	Score            int64     `json:"score" bson:"score"`
	OwnerUserId      uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	OwnerDisplayName string    `json:"ownerDisplayName" bson:"ownerDisplayName"`
	OwnerAvatar      string    `json:"ownerAvatar" bson:"ownerAvatar"`
	PostId           uuid.UUID `json:"postId" bson:"postId"`
	Text             string    `json:"text" bson:"text"`
	Deleted          bool      `json:"deleted" bson:"deleted"`
	DeletedDate      int64     `json:"deletedDate" bson:"deletedDate"`
	CreatedDate      int64     `json:"created_date" bson:"created_date"`
	LastUpdated      int64     `json:"last_updated" bson:"last_updated"`
}
