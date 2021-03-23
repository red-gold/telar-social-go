package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	domain "github.com/red-gold/ts-serverless/micros/vang/dto"
	models "github.com/red-gold/ts-serverless/micros/vang/models"
	service "github.com/red-gold/ts-serverless/micros/vang/services"
)

// UpdateMessageHandle handle create a new vang
func UpdateMessageHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.MessageModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		// Create service
		messageService, serviceErr := service.NewMessageService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		updatedMessage := &domain.Message{
			ObjectId:    model.ObjectId,
			OwnerUserId: req.UserID,
			RoomId:      model.RoomId,
			Text:        model.Text,
			CreatedDate: model.CreatedDate,
			UpdatedDate: model.UpdatedDate,
		}

		if err := messageService.UpdateMessageById(updatedMessage); err != nil {
			errorMessage := fmt.Sprintf("Update Message Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("updateMessageError", errorMessage)}, nil
		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
