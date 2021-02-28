package models

import uuid "github.com/satori/go.uuid"

type RelCirclesModel struct {
	CircleIds []string  `json:"circleIds"`
	RightId   uuid.UUID `json:"rightId"`
}
