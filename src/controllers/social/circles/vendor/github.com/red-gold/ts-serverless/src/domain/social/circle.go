package dto

import uuid "github.com/satori/go.uuid"

type Circle struct {
	ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
	CreatedDate int64     `json:"created_date" bson:"created_date"`
	OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	Name        string    `json:"name" bson:"name"`
	IsSystem    bool      `json:"isSystem" bson:"isSystem"`
}
