package service

import (
	uuid "github.com/gofrs/uuid"
	repo "github.com/red-gold/telar-core/data"
	dto "github.com/red-gold/ts-serverless/micros/gallery/dto"
)

type MediaService interface {
	SaveMedia(media *dto.Media) error
	SaveManyMedia(medias []dto.Media) error
	FindOneMedia(filter interface{}) (*dto.Media, error)
	FindMediaList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Media, error)
	QueryMedia(search string, ownerUserId *uuid.UUID, mediaTypeId *int, sortBy string, page int64) ([]dto.Media, error)
	FindById(objectId uuid.UUID) (*dto.Media, error)
	FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Media, error)
	UpdateMedia(filter interface{}, data interface{}, opts ...*repo.UpdateOptions) error
	UpdateMediaById(data *dto.Media) error
	DeleteMedia(filter interface{}) error
	DeleteMediaByOwner(ownerUserId uuid.UUID, mediaId uuid.UUID) error
	DeleteManyMedia(filter interface{}) error
	CreateMediaIndex(indexes map[string]interface{}) error
	FindByDirectory(ownerUserId uuid.UUID, directory string, limit int64, skip int64) ([]dto.Media, error)
	QueryAlbum(ownerUserId uuid.UUID, albumId *uuid.UUID, page int64, limit int64, sortBy string) ([]dto.Media, error)
	DeleteMediaByDirectory(ownerUserId uuid.UUID, directory string) error
}
