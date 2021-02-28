package models

import (
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/ts-serverless/constants"
)

type CreateMediaModel struct {
	ObjectId       uuid.UUID                     `json:"objectId"`
	DeletedDate    int64                         `json:"deletedDate"`
	CreatedDate    int64                         `json:"created_date"`
	Thumbnail      string                        `json:"thumbnail"`
	URL            string                        `json:"url"`
	FullPath       string                        `json:"fullPath"`
	Caption        string                        `json:"caption"`
	Directory      string                        `json:"directory"`
	FileName       string                        `json:"fileName"`
	OwnerUserId    uuid.UUID                     `json:"ownerUserId"`
	LastUpdated    int64                         `json:"last_updated"`
	AlbumId        uuid.UUID                     `json:"albumId"`
	Width          int64                         `json:"width"`
	Height         int64                         `json:"height"`
	Meta           string                        `json:"meta"`
	AccessUserList []string                      `json:"accessUserList"`
	Permission     constants.UserPermissionConst `json:"permission"`
	Deleted        bool                          `json:"deleted"`
}
