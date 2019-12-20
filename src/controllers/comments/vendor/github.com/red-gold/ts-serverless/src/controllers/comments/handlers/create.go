package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	notificationsModels "github.com/red-gold/telar-web/src/models/notifications"
	domain "github.com/red-gold/ts-serverless/src/domain/comments"
	models "github.com/red-gold/ts-serverless/src/models/comments"
	service "github.com/red-gold/ts-serverless/src/services/comments"
	uuid "github.com/satori/go.uuid"
)

type PostModelNotification struct {
	ObjectId         uuid.UUID `json:"objectId"`
	OwnerUserId      uuid.UUID `json:"ownerUserId"`
	OwnerDisplayName string    `json:"ownerDisplayName"`
	OwnerAvatar      string    `json:"ownerAvatar"`
}

// CreateCommentHandle handle create a new comment
func CreateCommentHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.CreateCommentModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal CreateCommentModel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelMarshalError", errorMessage)}, nil
		}

		if model.Text == "" {
			errorMessage := fmt.Sprintf("Comment text is required")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("commentTextEmptyError", errorMessage)}, nil
		}

		if model.PostId == uuid.Nil {
			errorMessage := fmt.Sprintf("Comment postId is required")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("commentpostIdNilError", errorMessage)}, nil
		}

		// Create service
		commentService, serviceErr := service.NewCommentService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("comment Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("commentServiceError", errorMessage)}, nil
		}

		newComment := &domain.Comment{
			OwnerUserId:      req.UserID,
			PostId:           model.PostId,
			Score:            0,
			Text:             model.Text,
			OwnerDisplayName: req.DisplayName,
			OwnerAvatar:      req.Avatar,
			Deleted:          false,
			DeletedDate:      0,
			CreatedDate:      utils.UTCNowUnix(),
			LastUpdated:      0,
		}

		if err := commentService.SaveComment(newComment); err != nil {
			errorMessage := fmt.Sprintf("Save Comment Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("saveCommentError", errorMessage)}, nil
		}

		// Create user headers for http request
		userHeaders := make(map[string][]string)
		userHeaders["uid"] = []string{req.UserID.String()}
		userHeaders["email"] = []string{req.Username}
		userHeaders["avatar"] = []string{req.Avatar}
		userHeaders["displayName"] = []string{req.DisplayName}
		userHeaders["role"] = []string{req.SystemRole}

		// Create request to increase comment counter on post
		go func() {

			postCommentURL := fmt.Sprintf("/posts/comment/+1/%s", model.PostId)
			_, postCommentErr := functionCall(http.MethodPut, []byte(""), postCommentURL, userHeaders)

			if postCommentErr != nil {
				messageError := fmt.Sprintf("Cannot save comment count on post! error: %s", postCommentErr.Error())
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
				// Should not send notification if the owner of the comment is same as owner of post
				return
			}
			URL := fmt.Sprintf("/%s/posts/%s", req.UserID, model.PostId)
			notificationModel := &notificationsModels.CreateNotificationModel{
				OwnerUserId:          req.UserID,
				OwnerDisplayName:     req.DisplayName,
				OwnerAvatar:          req.Avatar,
				Description:          fmt.Sprintf("%s commented on your post.", req.DisplayName),
				URL:                  URL,
				NotifyRecieverUserId: post.OwnerUserId,
				TargetId:             model.PostId,
				IsSeen:               false,
				Type:                 "comment",
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
			Body:       []byte(fmt.Sprintf(`{"success": true, "objectId": "%s"}`, newComment.ObjectId.String())),
			StatusCode: http.StatusOK,
		}, nil
	}
}
