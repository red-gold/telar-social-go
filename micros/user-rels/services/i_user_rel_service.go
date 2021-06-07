package service

import (
	uuid "github.com/gofrs/uuid"
	dto "github.com/red-gold/ts-serverless/micros/user-rels/dto"
)

type UserRelService interface {
	SaveUserRel(userRel *dto.UserRel) error
	FindOneUserRel(filter interface{}) (*dto.UserRel, error)
	FindUserRelList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.UserRel, error)
	QueryUserRel(search string, ownerUserId *uuid.UUID, sortBy string, page int64) ([]dto.UserRel, error)
	FindById(objectId uuid.UUID) (*dto.UserRel, error)
	FindByOwnerUserId(userId string) (*dto.UserRel, error)
	UpdateUserRel(filter interface{}, data interface{}) error
	UpdateUserRelById(data *dto.UserRel) error
	DeleteUserRel(filter interface{}) error
	DeleteUserRelByOwner(ownerUserId uuid.UUID, userRelId uuid.UUID) error
	DeleteManyUserRel(filter interface{}) error
	CreateUserRelIndex(indexes map[string]interface{}) error
	GetFollowers(userId uuid.UUID) ([]dto.UserRel, error)
	GetFollowing(userId uuid.UUID) ([]dto.UserRel, error)
	FollowUser(leftUser dto.UserRelMeta, rightUser dto.UserRelMeta, circleIds []string, tags []string) error
	UpdateRelCircles(leftId uuid.UUID, rightId uuid.UUID, circleIds []string) error
	UnfollowUser(leftId uuid.UUID, rightId uuid.UUID) error
	DeleteCircle(circleId string) error
}
