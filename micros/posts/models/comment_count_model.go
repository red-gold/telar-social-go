package models

import uuid "github.com/gofrs/uuid"

type CommentCountModel struct {
	PostId uuid.UUID `json:"postId"`
	Count  int       `json:"count"`
}
