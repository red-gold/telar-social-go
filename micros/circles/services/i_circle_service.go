package service

import (
	uuid "github.com/gofrs/uuid"
	dto "github.com/red-gold/ts-serverless/micros/circles/dto"
)

type CircleService interface {
	SaveCircle(circle *dto.Circle) error
	FindOneCircle(filter interface{}) (*dto.Circle, error)
	FindCircleList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Circle, error)
	QueryCircle(search string, ownerUserId *uuid.UUID, sortBy string, page int64) ([]dto.Circle, error)
	FindById(objectId uuid.UUID) (*dto.Circle, error)
	FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Circle, error)
	UpdateCircle(filter interface{}, data interface{}) error
	UpdateCircleById(data *dto.Circle) error
	DeleteCircle(filter interface{}) error
	DeleteCircleByOwner(ownerUserId uuid.UUID, circleId uuid.UUID) error
	DeleteManyCircle(filter interface{}) error
	CreateCircleIndex(indexes map[string]interface{}) error
}
