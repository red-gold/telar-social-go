package service

import (
	uuid "github.com/gofrs/uuid"
	coreData "github.com/red-gold/telar-core/data"
	dto "github.com/red-gold/ts-serverless/micros/votes/dto"
)

type VoteService interface {
	SaveVote(vote *dto.Vote) <-chan SaveResultAsync
	FindOneVote(filter interface{}) (*dto.Vote, error)
	FindVoteList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Vote, error)
	FindById(objectId uuid.UUID) (*dto.Vote, error)
	FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Vote, error)
	UpdateVote(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error
	UpdateVoteById(data *dto.Vote) error
	DeleteVote(filter interface{}) error
	DeleteVoteByOwner(ownerUserId uuid.UUID, voteId uuid.UUID) error
	DeleteManyVotes(filter interface{}) error
	CreateVoteIndex(indexes map[string]interface{}) error
	GetVoteByPostId(postId *uuid.UUID, sortBy string, page int64) ([]dto.Vote, error)
	DeleteVotesByPostId(ownerUserId uuid.UUID, postId uuid.UUID) error
}
