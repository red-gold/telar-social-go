package models

import uuid "github.com/gofrs/uuid"

type RelCirclesModel struct {
	CircleIds []string  `json:"circleIds"`
	RightId   uuid.UUID `json:"rightId"`
}
