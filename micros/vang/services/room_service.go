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

// RoomService handlers with injected dependencies
type RoomServiceImpl struct {
	RoomRepo repo.Repository
}

// NewRoomService initializes RoomService's dependencies and create new RoomService struct
func NewRoomService(db interface{}) (RoomService, error) {

	roomService := &RoomServiceImpl{}

	switch *config.AppConfig.DBType {
	case config.DB_MONGO:

		mongodb := db.(mongodb.MongoDatabase)
		roomService.RoomRepo = mongoRepo.NewDataRepositoryMongo(mongodb)

	}

	return roomService, nil
}

// SaveRoom save the room
func (s RoomServiceImpl) SaveRoom(room *dto.Room) error {

	if room.ObjectId == uuid.Nil {
		var uuidErr error
		room.ObjectId, uuidErr = uuid.NewV4()
		if uuidErr != nil {
			return uuidErr
		}
	}

	if room.CreatedDate == 0 {
		room.CreatedDate = utils.UTCNowUnix()
	}

	result := <-s.RoomRepo.Save(vangRoomCollectionName, room)

	return result.Error
}

// FindOneRoom get one room
func (s RoomServiceImpl) FindOneRoom(filter interface{}) (*dto.Room, error) {

	result := <-s.RoomRepo.FindOne(vangRoomCollectionName, filter)
	if result.Error() != nil {
		if result.Error() == repo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Error()
	}

	var roomResult dto.Room
	errDecode := result.Decode(&roomResult)
	if errDecode != nil {
		return nil, fmt.Errorf("Error docoding on dto.Room")
	}
	return &roomResult, nil
}

// FindRoomList get all rooms by filter
func (s RoomServiceImpl) FindRoomList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Room, error) {

	result := <-s.RoomRepo.Find(vangRoomCollectionName, filter, limit, skip, sort)
	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var roomList []dto.Room
	for result.Next() {
		var room dto.Room
		errDecode := result.Decode(&room)
		if errDecode != nil {
			return nil, fmt.Errorf("Error docoding on dto.Room")
		}
		roomList = append(roomList, room)
	}

	return roomList, nil
}

// FindByOwnerUserId find by owner user id
func (s RoomServiceImpl) FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Room, error) {
	sortMap := make(map[string]int)
	sortMap["createdDate"] = -1
	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		OwnerUserId: ownerUserId,
	}
	return s.FindRoomList(filter, 0, 0, sortMap)
}

// FindById find by room id
func (s RoomServiceImpl) FindById(objectId uuid.UUID) (*dto.Room, error) {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}
	return s.FindOneRoom(filter)
}

// UpdateRoom update the room
func (s RoomServiceImpl) UpdateRoom(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error {

	result := <-s.RoomRepo.Update(vangRoomCollectionName, filter, data, opts...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateRoom update the room
func (s RoomServiceImpl) UpdateRoomById(data *dto.Room) error {
	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: data.ObjectId,
	}

	updateOperator := coreData.UpdateOperator{
		Set: data,
	}
	err := s.UpdateRoom(filter, updateOperator)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRoom delete room by filter
func (s RoomServiceImpl) DeleteRoom(filter interface{}) error {

	result := <-s.RoomRepo.Delete(vangRoomCollectionName, filter, true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteRoom delete room by ownerUserId and roomId
func (s RoomServiceImpl) DeleteRoomByOwner(ownerUserId uuid.UUID, roomId uuid.UUID) error {

	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    roomId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteRoom(filter)
	if err != nil {
		return err
	}
	return nil
}

// DeleteManyRoom delete many room by filter
func (s RoomServiceImpl) DeleteManyRoom(filter interface{}) error {

	result := <-s.RoomRepo.Delete(vangRoomCollectionName, filter, false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateRoomIndex create index for room search.
func (s RoomServiceImpl) CreateRoomIndex(indexes map[string]interface{}) error {
	result := <-s.RoomRepo.CreateIndex(vangRoomCollectionName, indexes)
	return result
}

// GetPeerRoom get all room by room ID
func (s RoomServiceImpl) GetPeerRoom(roomId uuid.UUID, members []string, deactivePeerId uuid.UUID) (*dto.Room, error) {

	// filters
	include := make(map[string]interface{})
	include["$in"] = members

	filter := make(map[string]interface{})
	filter["members"] = include
	filter["objectId"] = roomId

	if deactivePeerId != uuid.Nil {
		inDeactiveUsers := make(map[string]interface{})
		inDeactiveUsers["$in"] = []string{deactivePeerId.String()}
		filter["deactiveUsers"] = inDeactiveUsers
	}

	return s.FindOneRoom(filter)
}

// DeleteRoomByRoomId delete room by room id
func (s RoomServiceImpl) DeleteRoomByRoomId(ownerUserId uuid.UUID, roomId uuid.UUID) error {

	filter := struct {
		PostId      uuid.UUID `json:"roomId" bson:"roomId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		PostId:      roomId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteManyRoom(filter)
	if err != nil {
		return err
	}
	return nil
}

// FindOneRoomByMembers find one room by members
func (s RoomServiceImpl) FindOneRoomByMembers(userIds []string, roomType int8) (*dto.Room, error) {

	include := make(map[string]interface{})
	include["$all"] = userIds

	filter := make(map[string]interface{})
	filter["members"] = include
	filter["type"] = roomType

	return s.FindOneRoom(filter)
}

// GetRoomsByUserId Get rooms by user ID
func (s RoomServiceImpl) GetRoomsByUserId(userId string, roomType int8) ([]dto.Room, error) {
	sortMap := make(map[string]int)
	sortMap["updatedDate"] = -1

	include := make(map[string]interface{})
	include["$in"] = []string{userId}

	nin := make(map[string]interface{})
	nin["$nin"] = []string{userId}

	filter := make(map[string]interface{})
	filter["members"] = include
	filter["type"] = roomType
	filter["deactiveUsers"] = nin

	return s.FindRoomList(filter, 0, 0, sortMap)
}

// IncreaseMemberCount Increase member count
func (s RoomServiceImpl) IncreaseMemberCount(roomId uuid.UUID, amount int64) error {

	increase := make(map[string]interface{})
	increase["memberCount"] = amount

	data := make(map[string]interface{})
	increase["$inc"] = increase

	filter := make(map[string]interface{})
	filter["objectId"] = roomId

	return s.UpdateRoom(filter, data)
}

// UpdateMessageMeta Update message meta
func (s RoomServiceImpl) UpdateMessageMeta(roomId uuid.UUID, amount, createdDate int64, text, ownerId string) error {

	increase := make(map[string]interface{})
	increase["messageCount"] = amount

	data := make(map[string]interface{})
	data["$inc"] = increase

	setData := make(map[string]interface{})
	lastMessage := make(map[string]interface{})
	lastMessage["text"] = text
	lastMessage["ownerId"] = ownerId
	lastMessage["createdDate"] = createdDate
	setData["lastMessage"] = lastMessage
	data["$set"] = setData

	filter := make(map[string]interface{})
	filter["objectId"] = roomId

	options := &coreData.UpdateOptions{}
	options.SetUpsert(true)
	return s.UpdateRoom(filter, data, options)
}

// UpdateMemberRead Increase member read count and read member date to now
func (s RoomServiceImpl) UpdateMemberRead(roomId uuid.UUID, userId uuid.UUID, amount, messageCreatedDate int64) error {

	readCountField := fmt.Sprintf("readCount.%s", userId.String())
	readDateField := fmt.Sprintf("readDate.%s", userId.String())

	setData := make(map[string]interface{})
	data := make(map[string]interface{})
	setData[readDateField] = messageCreatedDate
	setData[readCountField] = amount
	data["$set"] = setData
	filter := make(map[string]interface{})
	filter["objectId"] = roomId

	options := &coreData.UpdateOptions{}
	options.SetUpsert(true)
	return s.UpdateRoom(filter, data, options)
}

// DeactiveUserRoom set user delete a room
func (s RoomServiceImpl) DeactiveUserRoom(roomId uuid.UUID, userId uuid.UUID) error {

	// date to update
	push := make(map[string]interface{})
	push["deactiveUsers"] = userId.String()

	data := make(map[string]interface{})
	data["$push"] = push

	// filters
	include := make(map[string]interface{})
	include["$in"] = []string{userId.String()}

	nin := make(map[string]interface{})
	nin["$nin"] = []string{userId.String()}

	filter := make(map[string]interface{})
	filter["members"] = include
	filter["objectId"] = roomId
	filter["deactiveUsers"] = nin

	options := &coreData.UpdateOptions{}
	options.SetUpsert(true)
	return s.UpdateRoom(filter, data, options)
}

// ActiveAllPeerRoom active all peer room members
func (s RoomServiceImpl) ActiveAllPeerRoom(roomId uuid.UUID, members []string, deactivePeerId uuid.UUID) error {

	setData := make(map[string]interface{})
	setData["deactiveUsers"] = []string{}

	data := make(map[string]interface{})
	data["$set"] = setData

	// filters
	include := make(map[string]interface{})
	include["$in"] = members

	filter := make(map[string]interface{})
	filter["members"] = include
	filter["objectId"] = roomId

	if deactivePeerId != uuid.Nil {
		inDeactiveUsers := make(map[string]interface{})
		inDeactiveUsers["$in"] = []string{deactivePeerId.String()}
		filter["deactiveUsers"] = inDeactiveUsers
	}

	options := &coreData.UpdateOptions{}
	options.SetUpsert(true)
	return s.UpdateRoom(filter, data, options)
}

// GetActiveRoom active all peer room members
func (s RoomServiceImpl) GetActiveRoom(roomId uuid.UUID, members []string) (*dto.Room, error) {

	size := make(map[string]interface{})
	size["$size"] = 0

	// filters
	include := make(map[string]interface{})
	include["$in"] = members

	filter := make(map[string]interface{})
	filter["members"] = include
	filter["objectId"] = roomId
	filter["deactiveUsers"] = size

	return s.FindOneRoom(filter)
}
