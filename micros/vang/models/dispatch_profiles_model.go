package models

import "github.com/gofrs/uuid"

type DispatchProfilesModel struct {
	UserIds   []string  `json:"userIds" bson:"userIds"`
	ReqUserId uuid.UUID `json:"reqUserId" bson:"reqUserId"`
}
