package dto

import (
	uuid "github.com/gofrs/uuid"
)

type Room struct {
	ObjectId    uuid.UUID        `json:"objectId" bson:"objectId"`
	Members     []string         `json:"members" bson:"members"`
	Type        int8             `json:"type" bson:"type"` // {0: peer, 1: multiple}
	Seen        map[string]int64 `json:"seen" bson:"seen"` // {'userId1': last_seen_date_time, 'userId2': last_seen_date_time}
	CreatedDate int64            `json:"createdDate" bson:"createdDate"`
	UpdatedDate int64            `json:"updatedDate" bson:"updatedDate"`
}
