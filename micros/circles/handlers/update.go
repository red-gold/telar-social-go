package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	domain "github.com/red-gold/ts-serverless/micros/circles/dto"
	service "github.com/red-gold/ts-serverless/micros/circles/services"
)

// UpdateCircleHandle handle create a new circle
func UpdateCircleHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model domain.Circle
		if err := json.Unmarshal(req.Body, &model); err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		if model.Name == "" {
			errorMessage := fmt.Sprintf("Circle name can not be empty.")
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("circleNameEmptyError", errorMessage)}, nil
		}

		// Create service
		circleService, serviceErr := service.NewCircleService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		if err := circleService.UpdateCircleById(&model); err != nil {
			errorMessage := fmt.Sprintf("Update Circle Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("updateCircleError", errorMessage)}, nil
		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
