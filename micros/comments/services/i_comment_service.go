package service

import (
	uuid "github.com/gofrs/uuid"
	coreData "github.com/red-gold/telar-core/data"
	dto "github.com/red-gold/ts-serverless/micros/comments/dto"
)

type CommentService interface {
	SaveComment(comment *dto.Comment) error
	FindOneComment(filter interface{}) (*dto.Comment, error)
	FindCommentList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Comment, error)
	QueryComment(search string, ownerUserId *uuid.UUID, commentTypeId *int, sortBy string, page int64) ([]dto.Comment, error)
	QueryCommentIncludeProfile(search string, ownerUserId *uuid.UUID, commentTypeId *int, sortBy string, page int64) ([]dto.Comment, error)
	FindById(objectId uuid.UUID) (*dto.Comment, error)
	FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Comment, error)
	UpdateComment(filter interface{}, data interface{}) error
	UpdateManyComment(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error
	UpdateCommentById(data *dto.Comment) error
	DeleteComment(filter interface{}) error
	DeleteCommentByOwner(ownerUserId uuid.UUID, commentId uuid.UUID) error
	DeleteManyComments(filter interface{}) error
	CreateCommentIndex(indexes map[string]interface{}) error
	GetCommentByPostId(postId *uuid.UUID, sortBy string, page int64) ([]dto.Comment, error)
	DeleteCommentsByPostId(ownerUserId uuid.UUID, postId uuid.UUID) error
	UpdateCommentProfile(ownerUserId uuid.UUID, ownerDisplayName string, ownerAvatar string) error
}
