package service

import (
	uuid "github.com/gofrs/uuid"
	coreData "github.com/red-gold/telar-core/data"
	dto "github.com/red-gold/ts-serverless/micros/vang/dto"
)

type MessageService interface {
	SaveMessage(vang *dto.Message) error
	SaveManyMessages(messages []dto.Message) error
	FindOneMessage(filter interface{}) (*dto.Message, error)
	FindMessageList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Message, error)
	FindById(objectId uuid.UUID) (*dto.Message, error)
	FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Message, error)
	UpdateMessage(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error
	UpdateMessageById(data *dto.Message) error
	DeleteMessage(filter interface{}) error
	DeleteMessageByOwner(ownerUserId uuid.UUID, vangId uuid.UUID) error
	DeleteManyMessage(filter interface{}) error
	CreateMessageIndex(indexes map[string]interface{}) error
	GetMessageByRoomId(roomId *uuid.UUID, sortBy string, page int64, lteDate int64, gteDate int64) ([]dto.Message, error)
	DeleteMessageByRoomId(ownerUserId uuid.UUID, roomId uuid.UUID) error
}
