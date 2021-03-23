package handlers

import (
	"fmt"
	"net/http"

	uuid "github.com/gofrs/uuid"
	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	service "github.com/red-gold/ts-serverless/micros/vang/services"
)

// DeleteMessageHandle handle delete a Message
func DeleteMessageHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /message/id/:messageId
		messageId := req.GetParamByName("messageId")
		if messageId == "" {
			errorMessage := fmt.Sprintf("Message Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("messageIdRequired", errorMessage)}, nil
		}
		fmt.Printf("\n Message ID: %s", messageId)
		messageUUID, uuidErr := uuid.FromString(messageId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}
		fmt.Printf("\n Message UUID: %s", messageUUID)
		// Create service
		messageService, serviceErr := service.NewMessageService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("Message Service Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("messageServiceError", errorMessage)}, nil

		}

		if err := messageService.DeleteMessageByOwner(req.UserID, messageUUID); err != nil {
			errorMessage := fmt.Sprintf("Delete Message Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("deleteMessageError", errorMessage)}, nil

		}

		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
