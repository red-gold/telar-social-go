package models

import (
	uuid "github.com/satori/go.uuid"
)

type CreateNotificationModel struct {
	ObjectId             uuid.UUID `json:"objectId"`
	OwnerUserId          uuid.UUID `json:"ownerUserId"`
	OwnerDisplayName     string    `json:"ownerDisplayName"`
	OwnerAvatar          string    `json:"ownerAvatar"`
	CreatedDate          int64     `json:"created_date"`
	Description          string    `json:"description"`
	URL                  string    `json:"url"`
	NotifyRecieverUserId uuid.UUID `json:"notifyRecieverUserId"`
	TargetId             uuid.UUID `json:"targetId"`
	IsSeen               bool      `json:"isSeen"`
	Type                 string    `json:"type"`
	EmailNotification    int16     `json:"emailNotification"`
}
