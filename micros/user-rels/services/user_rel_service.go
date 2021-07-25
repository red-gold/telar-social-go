package service

import (
	"fmt"

	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/config"
	coreData "github.com/red-gold/telar-core/data"
	repo "github.com/red-gold/telar-core/data"
	"github.com/red-gold/telar-core/data/mongodb"
	mongoRepo "github.com/red-gold/telar-core/data/mongodb"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/utils"
	dto "github.com/red-gold/ts-serverless/micros/user-rels/dto"
)

// UserRelService handlers with injected dependencies
type UserRelServiceImpl struct {
	UserRelRepo repo.Repository
}

// NewUserRelService initializes UserRelService's dependencies and create new UserRelService struct
func NewUserRelService(db interface{}) (UserRelService, error) {

	userRelService := &UserRelServiceImpl{}

	switch *config.AppConfig.DBType {
	case config.DB_MONGO:

		mongodb := db.(mongodb.MongoDatabase)
		userRelService.UserRelRepo = mongoRepo.NewDataRepositoryMongo(mongodb)

	}

	return userRelService, nil
}

// SaveUserRel save the userRel
func (s UserRelServiceImpl) SaveUserRel(userRel *dto.UserRel) error {

	if userRel.ObjectId == uuid.Nil {
		var uuidErr error
		userRel.ObjectId, uuidErr = uuid.NewV4()
		if uuidErr != nil {
			return uuidErr
		}
	}

	if userRel.CreatedDate == 0 {
		userRel.CreatedDate = utils.UTCNowUnix()
	}

	result := <-s.UserRelRepo.Save(userRelCollectionName, userRel)

	return result.Error
}

// FindOneUserRel get one userRel
func (s UserRelServiceImpl) FindOneUserRel(filter interface{}) (*dto.UserRel, error) {

	result := <-s.UserRelRepo.FindOne(userRelCollectionName, filter)
	if result.Error() != nil {
		return nil, result.Error()
	}

	var userRelResult dto.UserRel
	errDecode := result.Decode(&userRelResult)
	if errDecode != nil {
		return nil, fmt.Errorf("Error docoding on dto.UserRel")
	}
	return &userRelResult, nil
}

// FindUserRelList get all userRels by filter
func (s UserRelServiceImpl) FindUserRelList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.UserRel, error) {

	result := <-s.UserRelRepo.Find(userRelCollectionName, filter, limit, skip, sort)
	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var userRelList []dto.UserRel
	for result.Next() {
		var userRel dto.UserRel
		errDecode := result.Decode(&userRel)
		if errDecode != nil {
			return nil, fmt.Errorf("Error docoding on dto.UserRel")
		}
		userRelList = append(userRelList, userRel)
	}

	return userRelList, nil
}

// FindRelsIncludeProfile get all user relations by filter including user profile entity
func (s UserRelServiceImpl) FindRelsIncludeProfile(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.UserRel, error) {
	var pipeline []interface{}

	matchOperator := make(map[string]interface{})
	matchOperator["$match"] = filter

	sortOperator := make(map[string]interface{})
	sortOperator["$sort"] = sort

	pipeline = append(pipeline, matchOperator, sortOperator)

	if skip > 0 {
		skipOperator := make(map[string]interface{})
		skipOperator["$skip"] = skip
		pipeline = append(pipeline, skipOperator)
	}

	if limit > 0 {
		limitOperator := make(map[string]interface{})
		limitOperator["$limit"] = limit
		pipeline = append(pipeline, limitOperator)
	}

	// Add left user pipeline
	lookupLeftUser := make(map[string]map[string]string)
	lookupLeftUser["$lookup"] = map[string]string{
		"localField":   "leftId",
		"from":         "userProfile",
		"foreignField": "objectId",
		"as":           "leftUser",
	}

	unwindLeftUser := make(map[string]interface{})
	unwindLeftUser["$unwind"] = "$leftUser"
	pipeline = append(pipeline, lookupLeftUser, unwindLeftUser)

	// Add right user pipeline
	lookupRightUser := make(map[string]map[string]string)
	lookupRightUser["$lookup"] = map[string]string{
		"localField":   "rightId",
		"from":         "userProfile",
		"foreignField": "objectId",
		"as":           "rightUser",
	}

	unwindRightUser := make(map[string]interface{})
	unwindRightUser["$unwind"] = "$rightUser"
	pipeline = append(pipeline, lookupRightUser, unwindRightUser)
	log.Info("pipeline %v", pipeline)

	projectOperator := make(map[string]interface{})
	project := make(map[string]interface{})

	// Add project operator
	project["objectId"] = 1
	project["created_date"] = 1
	project["leftId"] = 1
	project["rightId"] = 1
	project["rel"] = 1
	project["tags"] = 1
	project["circleIds"] = 1
	// left user
	project["left.userId"] = "$leftId"
	project["left.fullName"] = "$leftUser.fullName"
	project["left.instagramId"] = "$leftUser.instagramId"
	project["left.twitterId"] = "$leftUser.twitterId"
	project["left.linkedInId"] = "$leftUser.linkedInId"
	project["left.facebookId"] = "$leftUser.facebookId"
	project["left.socialName"] = "$leftUser.socialName"
	project["left.created_date"] = "$leftUser.created_date"
	project["left.banner"] = "$leftUser.banner"
	project["left.avatar"] = "$leftUser.avatar"
	// Right user
	project["right.userId"] = "$rightId"
	project["right.fullName"] = "$rightUser.fullName"
	project["right.instagramId"] = "$rightUser.instagramId"
	project["right.twitterId"] = "$rightUser.twitterId"
	project["right.linkedInId"] = "$rightUser.linkedInId"
	project["right.facebookId"] = "$rightUser.facebookId"
	project["right.socialName"] = "$rightUser.socialName"
	project["right.created_date"] = "$rightUser.created_date"
	project["right.banner"] = "$rightUser.banner"
	project["right.avatar"] = "$rightUser.avatar"

	projectOperator["$project"] = project

	pipeline = append(pipeline, projectOperator)

	result := <-s.UserRelRepo.Aggregate(userRelCollectionName, pipeline)

	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var postList []dto.UserRel
	for result.Next() {
		var post dto.UserRel
		errDecode := result.Decode(&post)
		if errDecode != nil {
			return nil, fmt.Errorf("Error docoding on dto.UserRel")
		}
		postList = append(postList, post)
	}

	return postList, nil
}

// QueryUserRel get all userRels by query
func (s UserRelServiceImpl) QueryUserRel(search string, ownerUserId *uuid.UUID, sortBy string, page int64) ([]dto.UserRel, error) {
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
	fmt.Println(filter)
	result, err := s.FindUserRelList(filter, limit, skip, sortMap)

	return result, err
}

// FindByOwnerUserId find by owner user id
func (s UserRelServiceImpl) FindByOwnerUserId(ownerUserId string) (*dto.UserRel, error) {

	filter := struct {
		OwnerUserId string `json:"ownerUserId" bson:"ownerUserId"`
	}{
		OwnerUserId: ownerUserId,
	}
	return s.FindOneUserRel(filter)
}

// FindById find by userRel id
func (s UserRelServiceImpl) FindById(objectId uuid.UUID) (*dto.UserRel, error) {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}
	return s.FindOneUserRel(filter)
}

// UpdateUserRel update the userRel
func (s UserRelServiceImpl) UpdateUserRel(filter interface{}, data interface{}) error {

	result := <-s.UserRelRepo.Update(userRelCollectionName, filter, data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateUserRel update the userRel
func (s UserRelServiceImpl) UpdateUserRelById(data *dto.UserRel) error {
	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: data.ObjectId,
	}
	updateOperator := coreData.UpdateOperator{
		Set: data,
	}
	err := s.UpdateUserRel(filter, updateOperator)
	return err
}

// DeleteUserRel delete userRel by filter
func (s UserRelServiceImpl) DeleteUserRel(filter interface{}) error {

	result := <-s.UserRelRepo.Delete(userRelCollectionName, filter, true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteUserRel delete userRel by ownerUserId and userRelId
func (s UserRelServiceImpl) DeleteUserRelByOwner(ownerUserId uuid.UUID, userRelId uuid.UUID) error {

	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    userRelId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteUserRel(filter)
	if err != nil {
		return err
	}
	return nil
}

// DeleteManyUserRel delete many userRel by filter
func (s UserRelServiceImpl) DeleteManyUserRel(filter interface{}) error {

	result := <-s.UserRelRepo.Delete(userRelCollectionName, filter, false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateUserRelIndex create index for userRel search.
func (s UserRelServiceImpl) CreateUserRelIndex(indexes map[string]interface{}) error {
	result := <-s.UserRelRepo.CreateIndex(userRelCollectionName, indexes)
	return result
}

// GetFollowers Get user followers by userId
func (s UserRelServiceImpl) GetFollowers(userId uuid.UUID) ([]dto.UserRel, error) {
	sortMap := make(map[string]int)
	sortMap["created_date"] = -1
	filter := struct {
		RightId uuid.UUID `json:"rightId" bson:"rightId"`
	}{
		RightId: userId,
	}
	return s.FindRelsIncludeProfile(filter, 0, 0, sortMap)
}

// GetFollowing Get user's following by userId
func (s UserRelServiceImpl) GetFollowing(userId uuid.UUID) ([]dto.UserRel, error) {
	sortMap := make(map[string]int)
	sortMap["created_date"] = -1
	filter := struct {
		LeftId uuid.UUID `json:"leftId" bson:"leftId"`
	}{
		LeftId: userId,
	}
	return s.FindRelsIncludeProfile(filter, 0, 0, sortMap)
}

// FollowUser create relation between two users
func (s UserRelServiceImpl) FollowUser(leftUser dto.UserRelMeta, rightUser dto.UserRelMeta, circleIds []string, tags []string) error {

	newUserRel := &dto.UserRel{
		Left:      leftUser,
		LeftId:    leftUser.UserId,
		Right:     rightUser,
		RightId:   rightUser.UserId,
		Rel:       []string{leftUser.UserId.String(), rightUser.UserId.String()},
		CircleIds: circleIds,
		Tags:      tags,
	}
	err := s.SaveUserRel(newUserRel)
	return err
}

// UpdateRelCircles update the user relation circle ids
func (s UserRelServiceImpl) UpdateRelCircles(leftId uuid.UUID, rightId uuid.UUID, circleIds []string) error {
	filter := struct {
		LeftId  uuid.UUID `json:"leftId" bson:"leftId"`
		RightId uuid.UUID `json:"rightId" bson:"rightId"`
	}{
		LeftId:  leftId,
		RightId: rightId,
	}
	updateOperator := coreData.UpdateOperator{
		Set: struct {
			CircleIds []string `json:"circleIds" bson:"circleIds"`
		}{
			CircleIds: circleIds,
		},
	}
	err := s.UpdateUserRel(filter, updateOperator)
	return err
}

// DeleteCircle delete the circle from user-rel
func (s UserRelServiceImpl) DeleteCircle(circleId string) error {
	filter := struct{}{}
	pullOperator := make(map[string]interface{})
	inOperator := make(map[string]interface{})
	inOperator["$in"] = []string{circleId}
	circleIds := make(map[string]interface{})
	circleIds["circleIds"] = inOperator
	pullOperator["$pull"] = circleIds
	err := s.UpdateUserRel(filter, pullOperator)
	return err
}

// UnfollowUser delete relation between two users by left and right userId
func (s UserRelServiceImpl) UnfollowUser(leftId uuid.UUID, rightId uuid.UUID) error {

	filter := struct {
		LeftId  uuid.UUID `json:"leftId" bson:"leftId"`
		RightId uuid.UUID `json:"rightId" bson:"rightId"`
	}{
		LeftId:  leftId,
		RightId: rightId,
	}
	err := s.DeleteUserRel(filter)
	if err != nil {
		return err
	}
	return nil
}
