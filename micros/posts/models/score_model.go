package models

import uuid "github.com/gofrs/uuid"

type ScoreModel struct {
	PostId uuid.UUID `json:"postId"`
	Count  int       `json:"count"`
}
