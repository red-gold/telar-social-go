package models

import "github.com/gofrs/uuid"

type GetUserRoomsModel struct {
	UserId uuid.UUID `json:"userId" bson:"userId"`
}
