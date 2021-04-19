package models

type GetProfilesModel struct {
	UserIds []string `json:"userIds" bson:"userIds"`
}
