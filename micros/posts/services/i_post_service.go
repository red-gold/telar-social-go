package service

import (
	uuid "github.com/gofrs/uuid"
	repo "github.com/red-gold/telar-core/data"
	dto "github.com/red-gold/ts-serverless/micros/posts/dto"
	"github.com/red-gold/ts-serverless/micros/posts/models"
)

type PostService interface {
	SavePost(post *dto.Post) error
	FindOnePost(filter interface{}) (*dto.Post, error)
	FindPostList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Post, error)
	FindPostsIncludeProfile(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Post, error)
	QueryPost(search string, ownerUserIds []uuid.UUID, postTypeId int, sortBy string, page int64) ([]dto.Post, error)
	QueryPostIncludeUser(search string, ownerUserIds []uuid.UUID, postTypeId int, sortBy string, page int64) ([]dto.Post, error)
	FindById(objectId uuid.UUID) (*dto.Post, error)
	FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Post, error)
	FindByURLKey(urlKey string) (*dto.Post, error)
	UpdatePost(filter interface{}, data interface{}, opts ...*repo.UpdateOptions) error
	UpdateManyPost(filter interface{}, data interface{}, opts ...*repo.UpdateOptions) error
	UpdatePostById(data *models.PostUpdateModel) error
	DeletePost(filter interface{}) error
	DeletePostByOwner(ownerUserId uuid.UUID, postId uuid.UUID) error
	DeleteManyPost(filter interface{}) error
	CreatePostIndex(indexes map[string]interface{}) error
	DisableCommnet(OwnerUserId uuid.UUID, objectId uuid.UUID, value bool) error
	DisableSharing(OwnerUserId uuid.UUID, objectId uuid.UUID, value bool) error
	IncrementScoreCount(objectId uuid.UUID, ownerUserId uuid.UUID, avatar string) error
	DecrementScoreCount(objectId uuid.UUID, ownerUserId uuid.UUID) error
	Increment(objectId uuid.UUID, field string, value int) error
	IncrementCommentCount(objectId uuid.UUID) error
	DecerementCommentCount(objectId uuid.UUID) error
	UpdatePostProfile(ownerUserId uuid.UUID, ownerDisplayName string, ownerAvatar string) error
	UpdatePostURLKey(postId uuid.UUID, urlKey string) error
}
