package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	notificationsModels "github.com/red-gold/telar-web/src/models/notifications"
	domain "github.com/red-gold/ts-serverless/src/domain/social"
	socialModels "github.com/red-gold/ts-serverless/src/models/social"
	service "github.com/red-gold/ts-serverless/src/services/social"
)

// CreateUserRelHandle handle create a new userRel
func CreateUserRelHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model domain.UserRel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal CreateUserRelModel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelMarshalError", errorMessage)}, nil
		}

		// Create service
		userRelService, serviceErr := service.NewUserRelService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("userRel Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("userRelServiceError", errorMessage)}, nil
		}

		if err := userRelService.SaveUserRel(&model); err != nil {
			errorMessage := fmt.Sprintf("Save UserRel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("saveUserRelError", errorMessage)}, nil
		}

		return handler.Response{
			Body:       []byte(fmt.Sprintf(`{"success": true, "objectId": "%s"}`, model.ObjectId.String())),
			StatusCode: http.StatusOK,
		}, nil
	}
}

//FollowHandle handle create a new userRel
func FollowHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model socialModels.FollowModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal CreateUserRelModel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelMarshalError", errorMessage)}, nil
		}

		// Create service
		userRelService, serviceErr := service.NewUserRelService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("userRel Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("userRelServiceError", errorMessage)}, nil
		}

		// Left User Meta
		leftUserMeta := domain.UserRelMeta{
			UserId:   req.UserID,
			FullName: req.DisplayName,
			Avatar:   req.Avatar,
		}

		// Right User Meta
		rightUserMeta := domain.UserRelMeta{
			UserId:   model.RightUser.UserId,
			FullName: model.RightUser.FullName,
			Avatar:   model.RightUser.Avatar,
		}

		if err := userRelService.FollowUser(leftUserMeta, rightUserMeta, model.CircleIds); err != nil {
			errorMessage := fmt.Sprintf("Save UserRel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("saveUserRelError", errorMessage)}, nil
		}

		// Create notification
		go func() {

			// Create user headers for http request
			userHeaders := make(map[string][]string)
			userHeaders["uid"] = []string{req.UserID.String()}
			userHeaders["email"] = []string{req.Username}
			userHeaders["avatar"] = []string{req.Avatar}
			userHeaders["displayName"] = []string{req.DisplayName}
			userHeaders["role"] = []string{req.SystemRole}

			URL := fmt.Sprintf("/%s", req.UserID)
			notificationModel := &notificationsModels.CreateNotificationModel{
				OwnerUserId:          req.UserID,
				OwnerDisplayName:     req.DisplayName,
				OwnerAvatar:          req.Avatar,
				Description:          fmt.Sprintf("%s is following you.", req.DisplayName),
				URL:                  URL,
				NotifyRecieverUserId: model.RightUser.UserId,
				TargetId:             model.RightUser.UserId,
				IsSeen:               false,
				Type:                 "follow",
			}
			notificationBytes, marshalErr := json.Marshal(notificationModel)
			if marshalErr != nil {
				fmt.Printf("Cannot marshal notification! error: %s", marshalErr.Error())

			}

			notificationURL := "/notifications"
			_, notificationIndexErr := functionCall(http.MethodPost, notificationBytes, notificationURL, userHeaders)
			if notificationIndexErr != nil {
				fmt.Printf("\nCannot save notification on follow user! error: %s\n", notificationIndexErr.Error())
			}

		}()

		return handler.Response{
			Body:       []byte(fmt.Sprintf(`{"success": true}`)),
			StatusCode: http.StatusOK,
		}, nil
	}
}
