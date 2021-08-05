package service

import (
	uuid "github.com/gofrs/uuid"
	coreData "github.com/red-gold/telar-core/data"
	dto "github.com/red-gold/ts-serverless/micros/vang/dto"
)

type RoomService interface {
	SaveRoom(vang *dto.Room) error
	FindOneRoom(filter interface{}) (*dto.Room, error)
	FindRoomList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Room, error)
	FindById(objectId uuid.UUID) (*dto.Room, error)
	FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Room, error)
	UpdateRoom(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error
	UpdateRoomById(data *dto.Room) error
	DeleteRoom(filter interface{}) error
	DeleteRoomByOwner(ownerUserId uuid.UUID, vangId uuid.UUID) error
	DeleteManyRoom(filter interface{}) error
	CreateRoomIndex(indexes map[string]interface{}) error
	GetPeerRoom(roomId uuid.UUID, members []string, DeactivePeerId uuid.UUID) (*dto.Room, error)
	DeleteRoomByRoomId(ownerUserId uuid.UUID, roomId uuid.UUID) error
	FindOneRoomByMembers(userIds []string, roomType int8) (*dto.Room, error)
	GetRoomsByUserId(userId string, roomType int8) ([]dto.Room, error)
	UpdateMessageMeta(roomId uuid.UUID, amount, createdDate int64, text, ownerId string) error
	UpdateMemberRead(roomId uuid.UUID, userId uuid.UUID, amount, messageCreatedDate int64) error
	DeactiveUserRoom(roomId uuid.UUID, userId uuid.UUID) error
	ActiveAllPeerRoom(roomId uuid.UUID, members []string, deactivePeerId uuid.UUID) error
	GetActiveRoom(roomId uuid.UUID, members []string) (*dto.Room, error)
}
