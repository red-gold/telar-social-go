package handlers

import (
	"fmt"
	"net/http"

	uuid "github.com/gofrs/uuid"
	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	service "github.com/red-gold/ts-serverless/micros/comments/services"
)

// DeleteCommentHandle handle delete a Comment
func DeleteCommentHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /comments/id/:commentId/post/:postId
		commentId := req.GetParamByName("commentId")
		if commentId == "" {
			errorMessage := fmt.Sprintf("Comment Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("commentIdRequired", errorMessage)}, nil
		}
		fmt.Printf("\n Comment ID: %s", commentId)
		commentUUID, uuidErr := uuid.FromString(commentId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}
		fmt.Printf("\n Comment UUID: %s", commentUUID)

		postId := req.GetParamByName("postId")
		if postId == "" {
			errorMessage := fmt.Sprintf("Post Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("postIdRequired", errorMessage)}, nil
		}
		// Create service
		commentService, serviceErr := service.NewCommentService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("Comment Service Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("commentServiceError", errorMessage)}, nil

		}

		if err := commentService.DeleteCommentByOwner(req.UserID, commentUUID); err != nil {
			errorMessage := fmt.Sprintf("Delete Comment Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("deleteCommentError", errorMessage)}, nil

		}

		// Create user headers for http request
		userHeaders := make(map[string][]string)
		userHeaders["uid"] = []string{req.UserID.String()}
		userHeaders["email"] = []string{req.Username}
		userHeaders["avatar"] = []string{req.Avatar}
		userHeaders["displayName"] = []string{req.DisplayName}
		userHeaders["role"] = []string{req.SystemRole}

		postIndexURL := fmt.Sprintf("/posts/comment/-1/%s", postId)
		_, postIndexErr := functionCall(http.MethodPut, []byte(""), postIndexURL, userHeaders)

		if postIndexErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError,
					Body: utils.MarshalError("decreasePostCommentCountError",
						fmt.Sprintf("Cannot save vote on post! error: %s", postIndexErr.Error()))},
				nil
		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// DeleteCommentByPostIdHandle handle delete a Comment but postId
func DeleteCommentByPostIdHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /Comments/post/:postId
		postId := req.GetParamByName("postId")
		if postId == "" {
			errorMessage := fmt.Sprintf("Post Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("postIdRequired", errorMessage)}, nil
		}
		PostUUID, uuidErr := uuid.FromString(postId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}

		// Create service
		commentService, serviceErr := service.NewCommentService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("Comment Service Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("CommentServiceError", errorMessage)}, nil

		}

		if err := commentService.DeleteCommentsByPostId(req.UserID, PostUUID); err != nil {
			errorMessage := fmt.Sprintf("Delete Comment Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("deleteCommentError", errorMessage)}, nil
		}

		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
