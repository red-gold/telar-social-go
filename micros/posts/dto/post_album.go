package dto

import uuid "github.com/gofrs/uuid"

type PostAlbum struct {
	Count   int       `json:"count" bson:"count"`
	Cover   string    `json:"cover" bson:"cover"`
	CoverId uuid.UUID `json:"coverId" bson:"coverId"`
	Photos  []string  `json:"photos" bson:"photos"`
	Title   string    `json:"title" bson:"title"`
}
