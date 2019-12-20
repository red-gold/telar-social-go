package models

import uuid "github.com/satori/go.uuid"

type RelMetaModel struct {
	UserId   uuid.UUID `json:"userId"`
	FullName string    `json:"fullName"`
	Avatar   string    `json:"avatar"`
}
