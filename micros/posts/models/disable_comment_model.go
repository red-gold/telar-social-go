package models

import uuid "github.com/gofrs/uuid"

type DisableCommentModel struct {
	PostId uuid.UUID `json:"postId"`
	Status bool      `json:"status"`
}
