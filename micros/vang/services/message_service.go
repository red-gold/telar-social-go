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
	dto "github.com/red-gold/ts-serverless/micros/vang/dto"
)

// MessageService handlers with injected dependencies
type MessageServiceImpl struct {
	MessageRepo repo.Repository
}

// NewMessageService initializes MessageService's dependencies and create new MessageService struct
func NewMessageService(db interface{}) (MessageService, error) {

	messageService := &MessageServiceImpl{}

	switch *config.AppConfig.DBType {
	case config.DB_MONGO:

		mongodb := db.(mongodb.MongoDatabase)
		messageService.MessageRepo = mongoRepo.NewDataRepositoryMongo(mongodb)

	}

	return messageService, nil
}

// SaveMessage save the message
func (s MessageServiceImpl) SaveMessage(message *dto.Message) error {

	if message.ObjectId == uuid.Nil {
		var uuidErr error
		message.ObjectId, uuidErr = uuid.NewV4()
		if uuidErr != nil {
			return uuidErr
		}
	}

	if message.CreatedDate == 0 {
		message.CreatedDate = utils.UTCNowUnix()
	}

	result := <-s.MessageRepo.Save(vangMessageCollectionName, message)

	return result.Error
}

// SaveManyMessages save many messages
func (s MessageServiceImpl) SaveManyMessages(messages []dto.Message) error {

	// https://github.com/golang/go/wiki/InterfaceSlice
	var interfaceSlice []interface{} = make([]interface{}, len(messages))
	for i, d := range messages {
		if d.ObjectId == uuid.Nil {
			var uuidErr error
			d.ObjectId, uuidErr = uuid.NewV4()
			if uuidErr != nil {
				return uuidErr
			}
		}

		if d.CreatedDate == 0 {
			d.CreatedDate = utils.UTCNowUnix()
		}
		interfaceSlice[i] = d
	}

	result := <-s.MessageRepo.SaveMany(vangMessageCollectionName, interfaceSlice)

	return result.Error
}

// FindOneMessage get one message
func (s MessageServiceImpl) FindOneMessage(filter interface{}) (*dto.Message, error) {

	result := <-s.MessageRepo.FindOne(vangMessageCollectionName, filter)
	if result.Error() != nil {
		if result.Error() == repo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Error()
	}

	var messageResult dto.Message
	errDecode := result.Decode(&messageResult)
	if errDecode != nil {
		return nil, fmt.Errorf("Error docoding on dto.Message")
	}
	return &messageResult, nil
}

// FindMessageList get all messages by filter
func (s MessageServiceImpl) FindMessageList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Message, error) {

	result := <-s.MessageRepo.Find(vangMessageCollectionName, filter, limit, skip, sort)
	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var messageList []dto.Message
	for result.Next() {
		var message dto.Message
		errDecode := result.Decode(&message)
		if errDecode != nil {
			return nil, fmt.Errorf("Error docoding on dto.Message")
		}
		messageList = append(messageList, message)
	}

	return messageList, nil
}

// FindByOwnerUserId find by owner user id
func (s MessageServiceImpl) FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Message, error) {
	sortMap := make(map[string]int)
	sortMap["createdDate"] = -1
	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		OwnerUserId: ownerUserId,
	}
	return s.FindMessageList(filter, 0, 0, sortMap)
}

// FindById find by message id
func (s MessageServiceImpl) FindById(objectId uuid.UUID) (*dto.Message, error) {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}
	return s.FindOneMessage(filter)
}

// UpdateMessage update the message
func (s MessageServiceImpl) UpdateMessage(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error {

	result := <-s.MessageRepo.Update(vangMessageCollectionName, filter, data, opts...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateMessage update the message
func (s MessageServiceImpl) UpdateMessageById(data *dto.Message) error {
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
	err := s.UpdateMessage(filter, updateOperator)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMessage delete message by filter
func (s MessageServiceImpl) DeleteMessage(filter interface{}) error {

	result := <-s.MessageRepo.Delete(vangMessageCollectionName, filter, true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteMessage delete message by ownerUserId and messageId
func (s MessageServiceImpl) DeleteMessageByOwner(ownerUserId uuid.UUID, messageId uuid.UUID) error {

	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    messageId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteMessage(filter)
	if err != nil {
		return err
	}
	return nil
}

// DeleteManyMessage delete many message by filter
func (s MessageServiceImpl) DeleteManyMessage(filter interface{}) error {

	result := <-s.MessageRepo.Delete(vangMessageCollectionName, filter, false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateMessageIndex create index for message search.
func (s MessageServiceImpl) CreateMessageIndex(indexes map[string]interface{}) error {
	result := <-s.MessageRepo.CreateIndex(vangMessageCollectionName, indexes)
	return result
}

// GetMessageByRoomId get all message by room ID
func (s MessageServiceImpl) GetMessageByRoomId(roomId *uuid.UUID, sortBy string, page int64, lteDate int64, gteDate int64) ([]dto.Message, error) {
	sortMap := make(map[string]int)
	sortMap[sortBy] = -1
	skip := numberOfItems * (page - 1)
	limit := numberOfItems

	filter := make(map[string]interface{})

	if roomId != nil {
		filter["roomId"] = *roomId
	}

	if lteDate > 0 {
		lessEqualDate := make(map[string]interface{})
		lessEqualDate["$lte"] = lteDate
		filter["createdDate"] = lessEqualDate
	}

	if gteDate > 0 {
		lessEqualDate := make(map[string]interface{})
		lessEqualDate["$gte"] = lteDate
		filter["createdDate"] = lessEqualDate
	}

	result, err := s.FindMessageList(filter, limit, skip, sortMap)

	return result, err
}

// DeleteMessageByRoomId delete message by room id
func (s MessageServiceImpl) DeleteMessageByRoomId(ownerUserId uuid.UUID, roomId uuid.UUID) error {

	filter := struct {
		PostId      uuid.UUID `json:"roomId" bson:"roomId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		PostId:      roomId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteManyMessage(filter)
	if err != nil {
		return err
	}
	return nil
}
