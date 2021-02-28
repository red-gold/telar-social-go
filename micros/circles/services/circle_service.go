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
	dto "github.com/red-gold/ts-serverless/micros/circles/dto"
)

// CircleService handlers with injected dependencies
type CircleServiceImpl struct {
	CircleRepo repo.Repository
}

// NewCircleService initializes CircleService's dependencies and create new CircleService struct
func NewCircleService(db interface{}) (CircleService, error) {

	circleService := &CircleServiceImpl{}

	switch *config.AppConfig.DBType {
	case config.DB_MONGO:

		mongodb := db.(mongodb.MongoDatabase)
		circleService.CircleRepo = mongoRepo.NewDataRepositoryMongo(mongodb)

	}

	return circleService, nil
}

// SaveCircle save the circle
func (s CircleServiceImpl) SaveCircle(circle *dto.Circle) error {

	if circle.ObjectId == uuid.Nil {
		var uuidErr error
		circle.ObjectId, uuidErr = uuid.NewV4()
		if uuidErr != nil {
			return uuidErr
		}
	}

	if circle.CreatedDate == 0 {
		circle.CreatedDate = utils.UTCNowUnix()
	}

	result := <-s.CircleRepo.Save(circleCollectionName, circle)

	return result.Error
}

// FindOneCircle get one circle
func (s CircleServiceImpl) FindOneCircle(filter interface{}) (*dto.Circle, error) {

	result := <-s.CircleRepo.FindOne(circleCollectionName, filter)
	if result.Error() != nil {
		return nil, result.Error()
	}

	var circleResult dto.Circle
	errDecode := result.Decode(&circleResult)
	if errDecode != nil {
		return nil, fmt.Errorf("Error docoding on dto.Circle")
	}
	return &circleResult, nil
}

// FindCircleList get all circles by filter
func (s CircleServiceImpl) FindCircleList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Circle, error) {

	result := <-s.CircleRepo.Find(circleCollectionName, filter, limit, skip, sort)
	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var circleList []dto.Circle
	for result.Next() {
		var circle dto.Circle
		errDecode := result.Decode(&circle)
		if errDecode != nil {
			return nil, fmt.Errorf("Error docoding on dto.Circle")
		}
		circleList = append(circleList, circle)
	}

	return circleList, nil
}

// QueryCircle get all circles by query
func (s CircleServiceImpl) QueryCircle(search string, ownerUserId *uuid.UUID, sortBy string, page int64) ([]dto.Circle, error) {
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
	result, err := s.FindCircleList(filter, limit, skip, sortMap)

	return result, err
}

// FindByOwnerUserId find by owner user id
func (s CircleServiceImpl) FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Circle, error) {
	sortMap := make(map[string]int)
	sortMap["created_date"] = -1
	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		OwnerUserId: ownerUserId,
	}
	return s.FindCircleList(filter, 0, 0, sortMap)
}

// FindById find by circle id
func (s CircleServiceImpl) FindById(objectId uuid.UUID) (*dto.Circle, error) {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}
	return s.FindOneCircle(filter)
}

// UpdateCircle update the circle
func (s CircleServiceImpl) UpdateCircle(filter interface{}, data interface{}) error {

	result := <-s.CircleRepo.Update(circleCollectionName, filter, data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateCircle update the circle
func (s CircleServiceImpl) UpdateCircleById(data *dto.Circle) error {
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
	err := s.UpdateCircle(filter, updateOperator)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCircle delete circle by filter
func (s CircleServiceImpl) DeleteCircle(filter interface{}) error {

	result := <-s.CircleRepo.Delete(circleCollectionName, filter, true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteCircle delete circle by ownerUserId and circleId
func (s CircleServiceImpl) DeleteCircleByOwner(ownerUserId uuid.UUID, circleId uuid.UUID) error {

	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    circleId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteCircle(filter)
	if err != nil {
		return err
	}
	return nil
}

// DeleteManyCircle delete many circle by filter
func (s CircleServiceImpl) DeleteManyCircle(filter interface{}) error {

	result := <-s.CircleRepo.Delete(circleCollectionName, filter, false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateCircleIndex create index for circle search.
func (s CircleServiceImpl) CreateCircleIndex(indexes map[string]interface{}) error {
	result := <-s.CircleRepo.CreateIndex(circleCollectionName, indexes)
	return result
}
