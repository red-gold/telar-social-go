package dto

import (
	uuid "github.com/gofrs/uuid"
)

type Room struct {
	ObjectId      uuid.UUID              `json:"objectId" bson:"objectId"`
	Members       []string               `json:"members" bson:"members"`
	Type          int8                   `json:"type" bson:"type"`                   // {0: peer, 1: multiple}
	ReadDate      map[string]int64       `json:"readDate" bson:"readDate"`           // {'userId1': last_seen_date_time, 'userId2': last_seen_date_time}
	ReadCount     map[string]int64       `json:"readCount" bson:"readCount"`         // {'userId1': read_count, 'userId2': read_count}
	ReadMessageId map[string]string      `json:"readMessageId" bson:"readMessageId"` // {'userId1': 'message_id_234', 'userId2': 'message_id_2323'}
	DeactiveUsers []string               `json:"deactiveUsers" bson:"deactiveUsers"` // ['userId1', 'userId2']
	LastMessage   map[string]interface{} `json:"lastMessage" bson:"lastMessage"`     // {'text': 'message_text', 'ownerId': 'userId'}
	MemberCount   int64                  `json:"memberCount" bson:"memberCount"`
	MessageCount  int64                  `json:"messageCount" bson:"messageCount"`
	CreatedDate   int64                  `json:"createdDate" bson:"createdDate"`
	UpdatedDate   int64                  `json:"updatedDate" bson:"updatedDate"`
}
