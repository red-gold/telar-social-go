package service

import (
	"fmt"

	"github.com/red-gold/telar-core/config"
	coreData "github.com/red-gold/telar-core/data"
	repo "github.com/red-gold/telar-core/data"
	"github.com/red-gold/telar-core/data/mongodb"
	mongoRepo "github.com/red-gold/telar-core/data/mongodb"
	"github.com/red-gold/telar-core/utils"
	dto "github.com/red-gold/ts-serverless/src/domain/comments"
	uuid "github.com/satori/go.uuid"
)

// CommentService handlers with injected dependencies
type CommentServiceImpl struct {
	CommentRepo repo.Repository
}

// NewCommentService initializes CommentService's dependencies and create new CommentService struct
func NewCommentService(db interface{}) (CommentService, error) {

	commentService := &CommentServiceImpl{}

	switch *config.AppConfig.DBType {
	case config.DB_MONGO:

		mongodb := db.(mongodb.MongoDatabase)
		commentService.CommentRepo = mongoRepo.NewDataRepositoryMongo(mongodb)

	}

	return commentService, nil
}

// SaveComment save the comment
func (s CommentServiceImpl) SaveComment(comment *dto.Comment) error {

	if comment.ObjectId == uuid.Nil {
		var uuidErr error
		comment.ObjectId, uuidErr = uuid.NewV4()
		if uuidErr != nil {
			return uuidErr
		}
	}

	if comment.CreatedDate == 0 {
		comment.CreatedDate = utils.UTCNowUnix()
	}

	result := <-s.CommentRepo.Save(commentCollectionName, comment)

	return result.Error
}

// FindOneComment get one comment
func (s CommentServiceImpl) FindOneComment(filter interface{}) (*dto.Comment, error) {

	result := <-s.CommentRepo.FindOne(commentCollectionName, filter)
	if result.Error() != nil {
		return nil, result.Error()
	}
	if result.NoResult() {
		return nil, nil
	}
	var commentResult dto.Comment
	errDecode := result.Decode(&commentResult)
	if errDecode != nil {
		return nil, fmt.Errorf("Error docoding on dto.Comment")
	}
	return &commentResult, nil
}

// FindCommentList get all comments by filter
func (s CommentServiceImpl) FindCommentList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Comment, error) {

	result := <-s.CommentRepo.Find(commentCollectionName, filter, limit, skip, sort)
	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var commentList []dto.Comment
	for result.Next() {
		var comment dto.Comment
		errDecode := result.Decode(&comment)
		if errDecode != nil {
			return nil, fmt.Errorf("Error docoding on dto.Comment")
		}
		commentList = append(commentList, comment)
	}

	return commentList, nil
}

// QueryComment get all comments by query
func (s CommentServiceImpl) QueryComment(search string, ownerUserId *uuid.UUID, commentTypeId *int, sortBy string, page int64) ([]dto.Comment, error) {
	sortMap := make(map[string]int)
	sortMap[sortBy] = -1
	skip := numberOfItems * (page - 1)
	limit := numberOfItems

	filter := make(map[string]interface{})
	if search != "" {
		filter["$text"] = coreData.SearchOperator{Search: search}
	}
	if ownerUserId != nil {
		filter["ownerUserId"] = *ownerUserId
	}
	if commentTypeId != nil {
		filter["commentTypeId"] = *commentTypeId
	}
	fmt.Println(filter)
	result, err := s.FindCommentList(filter, limit, skip, sortMap)

	return result, err
}

// FindByOwnerUserId find by owner user id
func (s CommentServiceImpl) FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Comment, error) {
	sortMap := make(map[string]int)
	sortMap["created_date"] = -1
	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		OwnerUserId: ownerUserId,
	}
	return s.FindCommentList(filter, 0, 0, sortMap)
}

// FindById find by comment id
func (s CommentServiceImpl) FindById(objectId uuid.UUID) (*dto.Comment, error) {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}
	return s.FindOneComment(filter)
}

// UpdateComment update the comment
func (s CommentServiceImpl) UpdateComment(filter interface{}, data interface{}) error {

	result := <-s.CommentRepo.Update(commentCollectionName, filter, data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateManyComment update the comment
func (s CommentServiceImpl) UpdateManyComment(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error {

	result := <-s.CommentRepo.UpdateMany(commentCollectionName, filter, data, opts...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateComment update the comment
func (s CommentServiceImpl) UpdateCommentById(data *dto.Comment) error {
	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    data.ObjectId,
		OwnerUserId: data.OwnerUserId,
	}
	data.LastUpdated = utils.UTCNowUnix()
	updateOperator := coreData.UpdateOperator{
		Set: data,
	}
	err := s.UpdateComment(filter, updateOperator)
	if err != nil {
		return err
	}
	return nil
}

// DeleteComment delete comment by filter
func (s CommentServiceImpl) DeleteComment(filter interface{}) error {

	result := <-s.CommentRepo.Delete(commentCollectionName, filter, true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteCommentByOwner delete comment by ownerUserId and commentId
func (s CommentServiceImpl) DeleteCommentByOwner(ownerUserId uuid.UUID, commentId uuid.UUID) error {

	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    commentId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteComment(filter)
	if err != nil {
		return err
	}
	return nil
}

// DeleteManyComment delete many comments by filter
func (s CommentServiceImpl) DeleteManyComments(filter interface{}) error {

	result := <-s.CommentRepo.Delete(commentCollectionName, filter, false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateCommentIndex create index for comment search.
func (s CommentServiceImpl) CreateCommentIndex(indexes map[string]interface{}) error {
	result := <-s.CommentRepo.CreateIndex(commentCollectionName, indexes)
	return result
}

// GetCommentByPostId get all comments by postId
func (s CommentServiceImpl) GetCommentByPostId(postId *uuid.UUID, sortBy string, page int64) ([]dto.Comment, error) {
	sortMap := make(map[string]int)
	sortMap[sortBy] = -1
	skip := numberOfItems * (page - 1)
	limit := numberOfItems

	filter := make(map[string]interface{})

	if postId != nil {
		filter["postId"] = *postId
	}

	result, err := s.FindCommentList(filter, limit, skip, sortMap)

	return result, err
}

// DeleteCommentsByPostId delete comments by postId
func (s CommentServiceImpl) DeleteCommentsByPostId(ownerUserId uuid.UUID, postId uuid.UUID) error {

	filter := struct {
		PostId      uuid.UUID `json:"postId" bson:"postId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		PostId:      postId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteManyComments(filter)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCommentProfile update the post
func (s CommentServiceImpl) UpdateCommentProfile(ownerUserId uuid.UUID, ownerDisplayName string, ownerAvatar string) error {
	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		OwnerUserId: ownerUserId,
	}

	data := struct {
		OwnerDisplayName string `json:"ownerDisplayName" bson:"ownerDisplayName"`
		OwnerAvatar      string `json:"ownerAvatar" bson:"ownerAvatar"`
	}{
		OwnerDisplayName: ownerDisplayName,
		OwnerAvatar:      ownerAvatar,
	}

	updateOperator := coreData.UpdateOperator{
		Set: data,
	}
	err := s.UpdateManyComment(filter, updateOperator)
	if err != nil {
		return err
	}
	return nil
}
