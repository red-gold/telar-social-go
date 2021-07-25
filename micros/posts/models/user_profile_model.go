package models

import (
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-web/constants"
)

type UserProfileModel struct {
	ObjectId       uuid.UUID                     `json:"objectId" bson:"objectId"`
	FullName       string                        `json:"fullName" bson:"fullName"`
	SocialName     string                        `json:"socialName" bson:"socialName"`
	Avatar         string                        `json:"avatar" bson:"avatar"`
	Banner         string                        `json:"banner" bson:"banner"`
	TagLine        string                        `json:"tagLine" bson:"tagLine"`
	CreatedDate    int64                         `json:"created_date" bson:"created_date"`
	LastUpdated    int64                         `json:"last_updated" bson:"last_updated"`
	Email          string                        `json:"email" bson:"email"`
	Birthday       int64                         `json:"birthday" bson:"birthday"`
	WebUrl         string                        `json:"webUrl" bson:"webUrl"`
	CompanyName    string                        `json:"companyName" bson:"companyName"`
	VoteCount      int64                         `json:"voteCount" bson:"voteCount"`
	ShareCount     int64                         `json:"shareCount" bson:"shareCount"`
	FollowCount    int64                         `json:"followCount" bson:"followCount"`
	FollowerCount  int64                         `json:"followerCount" bson:"followerCount"`
	PostCount      int64                         `json:"postCount" bson:"postCount"`
	FacebookId     string                        `json:"facebookId" bson:"facebookId"`
	InstagramId    string                        `json:"instagramId" bson:"instagramId"`
	TwitterId      string                        `json:"twitterId" bson:"twitterId"`
	AccessUserList []string                      `json:"accessUserList" bson:"accessUserList"`
	Permission     constants.UserPermissionConst `json:"permission" bson:"permission"`
}
