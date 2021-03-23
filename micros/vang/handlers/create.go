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
		var model []models.MessageModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal []models.MessageModel %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelMarshalError", errorMessage)}, nil
		}

		// Create service
		messageService, serviceErr := service.NewMessageService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("vang Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("messageServiceError", errorMessage)}, nil
		}

		var messages []dto.Message
		for _, v := range model {
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

		if err := messageService.SaveManyMessages(messages); err != nil {
			errorMessage := fmt.Sprintf("Save many messages %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("saveMessageError", errorMessage)}, nil
		}

		return handler.Response{
			Body:       []byte(fmt.Sprintf(`{"success": true}`)),
			StatusCode: http.StatusOK,
		}, nil
	}
}
