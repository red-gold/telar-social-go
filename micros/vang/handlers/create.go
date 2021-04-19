package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/vang/dto"
	models "github.com/red-gold/ts-serverless/micros/vang/models"
	service "github.com/red-gold/ts-serverless/micros/vang/services"
)

// CreateMessageHandle handle create a new vang
func SaveMessages(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.SaveMessagesModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal []models.MessageModel %s", err.Error())
			println(errorMessage)
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelMarshalError", errorMessage)}, nil
		}

		// Room service
		roomService, roomServiceErr := service.NewRoomService(db)
		if roomServiceErr != nil {
			errorMessage := fmt.Sprintf("vang room service %s", roomServiceErr.Error())
			println(errorMessage)
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("roomServiceError", errorMessage)}, nil
		}

		// Message service
		messageService, messageServiceErr := service.NewMessageService(db)
		if messageServiceErr != nil {
			errorMessage := fmt.Sprintf("vang message %s", messageServiceErr.Error())
			println(errorMessage)
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("messageServiceError", errorMessage)}, nil
		}

		// Map message model to DTO
		var messages []dto.Message
		for _, v := range model.Messages {
			newMessage := dto.Message{
				ObjectId:    v.ObjectId,
				OwnerUserId: v.OwnerUserId,
				RoomId:      v.RoomId,
				Text:        v.Text,
				CreatedDate: utils.UTCNowUnix(),
				UpdatedDate: utils.UTCNowUnix(),
			}
			messages = append(messages, newMessage)
		}

		var maxDate int64
		var lastMessage *dto.Message
		// Get last messsage
		for _, v := range messages {
			if v.CreatedDate > maxDate {
				maxDate = v.CreatedDate
				lastMessage = &v
			}
		}

		// Increase room message count
		go func() {
			err := roomService.UpdateMessageMeta(model.RoomId, int64(len(messages)), lastMessage.CreatedDate, lastMessage.Text, req.UserID.String())
			if err != nil {
				errorMessage := fmt.Sprintf("vang IncreaseMessageCount %s", err.Error())
				println(errorMessage)
			}
		}()

		if err := messageService.SaveManyMessages(messages); err != nil {
			errorMessage := fmt.Sprintf("Save many messages %s", err.Error())
			println(errorMessage)
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("saveMessageError", errorMessage)}, nil
		}

		return handler.Response{
			Body:       []byte(fmt.Sprintf(`{"success": true}`)),
			StatusCode: http.StatusOK,
		}, nil
	}
}
