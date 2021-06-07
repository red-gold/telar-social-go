package dto

import uuid "github.com/gofrs/uuid"

type UserRel struct {
	ObjectId    uuid.UUID   `json:"objectId" bson:"objectId"`
	CreatedDate int64       `json:"created_date" bson:"created_date"`
	Left        UserRelMeta `json:"left" bson:"left"`
	LeftId      uuid.UUID   `json:"leftId" bson:"leftId"`
	Right       UserRelMeta `json:"right" bson:"right"`
	RightId     uuid.UUID   `json:"rightId" bson:"rightId"`
	Rel         []string    `json:"rel" bson:"rel"`
	Tags        []string    `json:"tags" bson:"tags"`
	CircleIds   []string    `json:"circleIds" bson:"circleIds"`
}
