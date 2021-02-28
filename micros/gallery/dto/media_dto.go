package dto

import (
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/ts-serverless/constants"
)

type Media struct {
	ObjectId       uuid.UUID                     `json:"objectId" bson:"objectId"`
	DeletedDate    int64                         `json:"deletedDate" bson:"deletedDate"`
	CreatedDate    int64                         `json:"created_date" bson:"created_date"`
	Thumbnail      string                        `json:"thumbnail" bson:"thumbnail"`
	URL            string                        `json:"url" bson:"url"`
	FullPath       string                        `json:"fullPath" bson:"fullPath"`
	Caption        string                        `json:"caption" bson:"caption"`
	Directory      string                        `json:"directory" bson:"directory"`
	FileName       string                        `json:"fileName" bson:"fileName"`
	OwnerUserId    uuid.UUID                     `json:"ownerUserId" bson:"ownerUserId"`
	LastUpdated    int64                         `json:"last_updated" bson:"last_updated"`
	AlbumId        uuid.UUID                     `json:"albumId" bson:"albumId"`
	Width          int64                         `json:"width" bson:"width"`
	Height         int64                         `json:"height" bson:"height"`
	Meta           string                        `json:"meta" bson:"meta"`
	AccessUserList []string                      `json:"accessUserList" bson:"accessUserList"`
	Permission     constants.UserPermissionConst `json:"permission" bson:"permission"`
	Deleted        bool                          `json:"deleted" bson:"deleted"`
}
