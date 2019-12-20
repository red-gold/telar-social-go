package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	notificationsModels "github.com/red-gold/telar-web/src/models/notifications"
	domain "github.com/red-gold/ts-serverless/src/domain/votes"
	models "github.com/red-gold/ts-serverless/src/models/votes"
	service "github.com/red-gold/ts-serverless/src/services/votes"
	uuid "github.com/satori/go.uuid"
)

type PostModelNotification struct {
	ObjectId         uuid.UUID `json:"objectId"`
	OwnerUserId      uuid.UUID `json:"ownerUserId"`
	OwnerDisplayName string    `json:"ownerDisplayName"`
	OwnerAvatar      string    `json:"ownerAvatar"`
}

// CreateVoteHandle handle create a new vote
func CreateVoteHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.CreateVoteModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal CreateVoteModel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelMarshalError", errorMessage)}, nil
		}

		// Create service
		voteService, serviceErr := service.NewVoteService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("vote Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("voteServiceError", errorMessage)}, nil
		}

		newVote := &domain.Vote{
			OwnerUserId:      req.UserID,
			PostId:           model.PostId,
			OwnerDisplayName: req.DisplayName,
			OwnerAvatar:      req.Avatar,
			TypeId:           model.TypeId,
		}

		if err := voteService.SaveVote(newVote); err != nil {
			errorMessage := fmt.Sprintf("Save Vote Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("saveVoteError", errorMessage)}, nil
		}

		// Create user headers for http request
		userHeaders := make(map[string][]string)
		userHeaders["uid"] = []string{req.UserID.String()}
		userHeaders["email"] = []string{req.Username}
		userHeaders["avatar"] = []string{req.Avatar}
		userHeaders["displayName"] = []string{req.DisplayName}
		userHeaders["role"] = []string{req.SystemRole}

		go func() {
			postURL := fmt.Sprintf("/posts/score/+1/%s", model.PostId)
			_, postErr := functionCall(http.MethodPut, []byte(""), postURL, userHeaders)

			if postErr != nil {
				messageError := fmt.Sprintf("Cannot save vote on post! error: %s", postErr.Error())
				fmt.Println(messageError)
			}

		}()

		// Create notification request
		go func() {
			postURL := fmt.Sprintf("/posts/%s", model.PostId)
			postBody, postErr := functionCall(http.MethodGet, []byte(""), postURL, userHeaders)

			if postErr != nil {
				messageError := fmt.Sprintf("Cannot get the post! error: %s", postErr.Error())
				fmt.Println(messageError)
			}

			var post PostModelNotification
			marshalErr := json.Unmarshal(postBody, &post)
			if marshalErr != nil {
				messageError := fmt.Sprintf("Cannot unmarshal the post! error: %s", marshalErr.Error())
				fmt.Println(messageError)
			}

			if post.OwnerUserId == req.UserID {
				// Should not send notification if the owner of the vote is same as owner of post
				return
			}

			URL := fmt.Sprintf("/%s/posts/%s", req.UserID, model.PostId)
			notificationModel := &notificationsModels.CreateNotificationModel{
				OwnerUserId:          req.UserID,
				OwnerDisplayName:     req.DisplayName,
				OwnerAvatar:          req.Avatar,
				Description:          fmt.Sprintf("%s like your post.", req.DisplayName),
				URL:                  URL,
				NotifyRecieverUserId: post.OwnerUserId,
				TargetId:             model.PostId,
				IsSeen:               false,
				Type:                 "like",
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
			Body:       []byte(fmt.Sprintf(`{"success": true, "objectId": "%s"}`, newVote.ObjectId.String())),
			StatusCode: http.StatusOK,
		}, nil
	}
}
