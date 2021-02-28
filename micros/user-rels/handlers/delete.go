package handlers

import (
	"fmt"
	"net/http"

	uuid "github.com/gofrs/uuid"
	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	service "github.com/red-gold/ts-serverless/micros/user-rels/services"
)

// DeleteUserRelHandle handle delete a userRel
func DeleteUserRelHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /user-rels/:userRelId
		userRelId := req.GetParamByName("userRelId")
		if userRelId == "" {
			errorMessage := fmt.Sprintf("UserRel Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("userRelIdRequired", errorMessage)}, nil
		}
		fmt.Printf("\n UserRel ID: %s", userRelId)
		userRelUUID, uuidErr := uuid.FromString(userRelId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}
		fmt.Printf("\n UserRel UUID: %s", userRelUUID)
		// Create service
		userRelService, serviceErr := service.NewUserRelService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("UserRel Service Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("userRelServiceError", errorMessage)}, nil

		}

		if err := userRelService.DeleteUserRelByOwner(req.UserID, userRelUUID); err != nil {
			errorMessage := fmt.Sprintf("Delete UserRel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("deleteUserRelError", errorMessage)}, nil

		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// UnfollowHandle handle delete a userRel
func UnfollowHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /user-rels/unfollow/:userId
		userFollowingId := req.GetParamByName("userId")
		if userFollowingId == "" {
			errorMessage := fmt.Sprintf("User Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("userIdRequired", errorMessage)}, nil
		}

		userFollowingUUID, uuidErr := uuid.FromString(userFollowingId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("userFollowingUUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("userFollowingUUIDError", errorMessage)}, nil
		}

		// Create service
		userRelService, serviceErr := service.NewUserRelService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("UserRel Service Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("userRelServiceError", errorMessage)}, nil

		}

		if err := userRelService.UnfollowUser(req.UserID, userFollowingUUID); err != nil {
			errorMessage := fmt.Sprintf("Delete UserRel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("deleteUserRelError", errorMessage)}, nil
		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// DeleteCircle handle delete a userRel
func DeleteCircle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /user-rels/circle/:circleId
		circleId := req.GetParamByName("circleId")
		if circleId == "" {
			errorMessage := fmt.Sprintf("Circle Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("circleId", errorMessage)}, nil
		}

		// Create service
		userRelService, serviceErr := service.NewUserRelService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("UserRel Service Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("userRelServiceError", errorMessage)}, nil

		}

		if err := userRelService.DeleteCircle(circleId); err != nil {
			errorMessage := fmt.Sprintf("Delete circle from user-rel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("deleteCircleUserRelError", errorMessage)}, nil
		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
