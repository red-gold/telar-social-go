package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	domain "github.com/red-gold/ts-serverless/src/domain/comments"
	models "github.com/red-gold/ts-serverless/src/models/comments"
	service "github.com/red-gold/ts-serverless/src/services/comments"
)

// UpdateCommentHandle handle create a new comment
func UpdateCommentHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.CommentModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		// Create service
		commentService, serviceErr := service.NewCommentService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		updatedComment := &domain.Comment{
			ObjectId:         model.ObjectId,
			OwnerUserId:      req.UserID,
			PostId:           model.PostId,
			Score:            model.Score,
			Text:             model.Text,
			OwnerDisplayName: req.DisplayName,
			OwnerAvatar:      req.Avatar,
			Deleted:          model.Deleted,
			DeletedDate:      model.DeletedDate,
			CreatedDate:      model.CreatedDate,
			LastUpdated:      model.LastUpdated,
		}

		if err := commentService.UpdateCommentById(updatedComment); err != nil {
			errorMessage := fmt.Sprintf("Update Comment Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("updateCommentError", errorMessage)}, nil
		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// UpdateCommentProfileHandle handle create a new post
func UpdateCommentProfileHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create service
		postService, serviceErr := service.NewCommentService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		postService.UpdateCommentProfile(req.UserID, req.DisplayName, req.Avatar)
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
