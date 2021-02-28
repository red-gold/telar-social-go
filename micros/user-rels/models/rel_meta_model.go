package models

import uuid "github.com/gofrs/uuid"

type RelMetaModel struct {
	UserId   uuid.UUID `json:"userId"`
	FullName string    `json:"fullName"`
	Avatar   string    `json:"avatar"`
}
