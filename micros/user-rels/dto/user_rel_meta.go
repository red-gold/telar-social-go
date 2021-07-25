package dto

import uuid "github.com/gofrs/uuid"

type UserRelMeta struct {
	UserId      uuid.UUID `json:"userId" bson:"userId"`
	CreatedDate int64     `json:"created_date" bson:"created_date"`
	FullName    string    `json:"fullName" bson:"fullName"`
	SocialName  string    `json:"socialName" bson:"socialName"`
	InstagramId string    `json:"instagramId" bson:"instagramId"`
	TwitterId   string    `json:"twitterId" bson:"twitterId"`
	FacebookId  string    `json:"facebookId" bson:"facebookId"`
	LinkedinId  string    `json:"linkedInId" bson:"linkedInId"`
	Banner      string    `json:"banner" bson:"banner"`
	Avatar      string    `json:"avatar" bson:"avatar"`
}
