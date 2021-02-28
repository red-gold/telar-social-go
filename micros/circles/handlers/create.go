package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	uuid "github.com/gofrs/uuid"
	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	domain "github.com/red-gold/ts-serverless/micros/circles/dto"
	models "github.com/red-gold/ts-serverless/micros/circles/models"
	service "github.com/red-gold/ts-serverless/micros/circles/services"
)

const followingCircleName = "Following"

// CreateCircleHandle handle create a new circle
func CreateCircleHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.CreateCircleModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal CreateCircleModel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelMarshalError", errorMessage)}, nil
		}

		if model.Name == followingCircleName {
			errorMessage := fmt.Sprintf("Can not user 'Following' as a circle name")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("followingCircleNameError", errorMessage)}, nil
		}

		// Create a new circle
		newCircle := &domain.Circle{
			OwnerUserId: req.UserID,
			Name:        model.Name,
			IsSystem:    false,
		}
		// Create service
		circleService, serviceErr := service.NewCircleService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("circle Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("circleServiceError", errorMessage)}, nil
		}

		if err := circleService.SaveCircle(newCircle); err != nil {
			errorMessage := fmt.Sprintf("Save Circle Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("saveCircleError", errorMessage)}, nil
		}

		return handler.Response{
			Body:       []byte(fmt.Sprintf(`{"success": true, "objectId": "%s"}`, newCircle.ObjectId.String())),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// CreateFollowingHandle handle create a new circle
func CreateFollowingHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// params from /circles/following/:userId
		userId := req.GetParamByName("userId")
		if userId == "" {
			errorMessage := fmt.Sprintf("User Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("userIdRequired", errorMessage)}, nil
		}
		fmt.Printf("\n Post ID: %s", userId)
		userUUID, uuidErr := uuid.FromString(userId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}
		// Create a new circle
		newCircle := &domain.Circle{
			OwnerUserId: userUUID,
			Name:        followingCircleName,
			IsSystem:    true,
		}

		// Create service
		circleService, serviceErr := service.NewCircleService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("circle Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("circleServiceError", errorMessage)}, nil
		}

		if err := circleService.SaveCircle(newCircle); err != nil {
			errorMessage := fmt.Sprintf("Save Circle Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("saveCircleError", errorMessage)}, nil
		}

		return handler.Response{
			Body:       []byte(fmt.Sprintf(`{"success": true, "objectId": "%s"}`, newCircle.ObjectId.String())),
			StatusCode: http.StatusOK,
		}, nil
	}
}
