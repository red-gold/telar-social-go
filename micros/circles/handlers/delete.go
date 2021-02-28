package handlers

import (
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	service "github.com/red-gold/ts-serverless/micros/circles/services"
	uuid "github.com/satori/go.uuid"
)

// DeleteCircleHandle handle delete a circle
func DeleteCircleHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /circles/:circleId
		circleId := req.GetParamByName("circleId")
		if circleId == "" {
			errorMessage := fmt.Sprintf("Circle Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("circleIdRequired", errorMessage)}, nil
		}
		fmt.Printf("\n Circle ID: %s", circleId)
		circleUUID, uuidErr := uuid.FromString(circleId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}
		fmt.Printf("\n Circle UUID: %s", circleUUID)
		// Create service
		circleService, serviceErr := service.NewCircleService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("Circle Service Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("circleServiceError", errorMessage)}, nil

		}

		if err := circleService.DeleteCircleByOwner(req.UserID, circleUUID); err != nil {
			errorMessage := fmt.Sprintf("Delete Circle Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("deleteCircleError", errorMessage)}, nil

		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
