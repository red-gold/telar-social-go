package service

import (
	"fmt"

	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/config"
	coreData "github.com/red-gold/telar-core/data"
	repo "github.com/red-gold/telar-core/data"
	"github.com/red-gold/telar-core/data/mongodb"
	mongoRepo "github.com/red-gold/telar-core/data/mongodb"
	"github.com/red-gold/telar-core/utils"
	dto "github.com/red-gold/ts-serverless/micros/votes/dto"
)

// VoteService handlers with injected dependencies
type VoteServiceImpl struct {
	VoteRepo repo.Repository
}

// NewVoteService initializes VoteService's dependencies and create new VoteService struct
func NewVoteService(db interface{}) (VoteService, error) {

	voteService := &VoteServiceImpl{}

	switch *config.AppConfig.DBType {
	case config.DB_MONGO:

		mongodb := db.(mongodb.MongoDatabase)
		voteService.VoteRepo = mongoRepo.NewDataRepositoryMongo(mongodb)

	}

	return voteService, nil
}

// SaveVote save the vote
func (s VoteServiceImpl) SaveVote(vote *dto.Vote) error {

	if vote.ObjectId == uuid.Nil {
		var uuidErr error
		vote.ObjectId, uuidErr = uuid.NewV4()
		if uuidErr != nil {
			return uuidErr
		}
	}

	if vote.CreatedDate == 0 {
		vote.CreatedDate = utils.UTCNowUnix()
	}

	result := <-s.VoteRepo.Save(voteCollectionName, vote)

	return result.Error
}

// FindOneVote get one vote
func (s VoteServiceImpl) FindOneVote(filter interface{}) (*dto.Vote, error) {

	result := <-s.VoteRepo.FindOne(voteCollectionName, filter)
	if result.Error() != nil {
		return nil, result.Error()
	}

	var voteResult dto.Vote
	errDecode := result.Decode(&voteResult)
	if errDecode != nil {
		return nil, fmt.Errorf("Error docoding on dto.Vote")
	}
	return &voteResult, nil
}

// FindVoteList get all votes by filter
func (s VoteServiceImpl) FindVoteList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Vote, error) {

	result := <-s.VoteRepo.Find(voteCollectionName, filter, limit, skip, sort)
	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var voteList []dto.Vote
	for result.Next() {
		var vote dto.Vote
		errDecode := result.Decode(&vote)
		if errDecode != nil {
			return nil, fmt.Errorf("Error docoding on dto.Vote")
		}
		voteList = append(voteList, vote)
	}

	return voteList, nil
}

// FindByOwnerUserId find by owner user id
func (s VoteServiceImpl) FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Vote, error) {
	sortMap := make(map[string]int)
	sortMap["created_date"] = -1
	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		OwnerUserId: ownerUserId,
	}
	return s.FindVoteList(filter, 0, 0, sortMap)
}

// FindById find by vote id
func (s VoteServiceImpl) FindById(objectId uuid.UUID) (*dto.Vote, error) {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}
	return s.FindOneVote(filter)
}

// UpdateVote update the vote
func (s VoteServiceImpl) UpdateVote(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error {

	result := <-s.VoteRepo.Update(voteCollectionName, filter, data, opts...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateVote update the vote
func (s VoteServiceImpl) UpdateVoteById(data *dto.Vote) error {
	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    data.ObjectId,
		OwnerUserId: data.OwnerUserId,
	}

	updateOperator := coreData.UpdateOperator{
		Set: data,
	}
	err := s.UpdateVote(filter, updateOperator)
	if err != nil {
		return err
	}
	return nil
}

// DeleteVote delete vote by filter
func (s VoteServiceImpl) DeleteVote(filter interface{}) error {

	result := <-s.VoteRepo.Delete(voteCollectionName, filter, true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteVote delete vote by ownerUserId and voteId
func (s VoteServiceImpl) DeleteVoteByOwner(ownerUserId uuid.UUID, voteId uuid.UUID) error {

	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    voteId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteVote(filter)
	if err != nil {
		return err
	}
	return nil
}

// DeleteManyVotes delete many votes by filter
func (s VoteServiceImpl) DeleteManyVotes(filter interface{}) error {

	result := <-s.VoteRepo.Delete(voteCollectionName, filter, false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateVoteIndex create index for vote search.
func (s VoteServiceImpl) CreateVoteIndex(indexes map[string]interface{}) error {
	result := <-s.VoteRepo.CreateIndex(voteCollectionName, indexes)
	return result
}

// GetVoteByPostId get all votes by postId
func (s VoteServiceImpl) GetVoteByPostId(postId *uuid.UUID, sortBy string, page int64) ([]dto.Vote, error) {
	sortMap := make(map[string]int)
	sortMap[sortBy] = -1
	skip := numberOfItems * (page - 1)
	limit := numberOfItems

	filter := make(map[string]interface{})

	if postId != nil {
		filter["postId"] = *postId
	}

	result, err := s.FindVoteList(filter, limit, skip, sortMap)

	return result, err
}

// DeleteVotesByPostId delete votes by postId
func (s VoteServiceImpl) DeleteVotesByPostId(ownerUserId uuid.UUID, postId uuid.UUID) error {

	filter := struct {
		PostId      uuid.UUID `json:"postId" bson:"postId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		PostId:      postId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteManyVotes(filter)
	if err != nil {
		return err
	}
	return nil
}
