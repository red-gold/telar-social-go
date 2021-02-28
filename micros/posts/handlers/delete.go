package handlers

import (
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	service "github.com/red-gold/ts-serverless/micros/posts/services"
	uuid "github.com/satori/go.uuid"
)

// DeletePostHandle handle delete a post
func DeletePostHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /posts/:postId
		postId := req.GetParamByName("postId")
		if postId == "" {
			errorMessage := fmt.Sprintf("Post Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("postIdRequired", errorMessage)}, nil
		}

		fmt.Printf("\n Post ID: %s", postId)
		postUUID, uuidErr := uuid.FromString(postId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}

		fmt.Printf("\n Post UUID: %s", postUUID)
		// Create service
		postService, serviceErr := service.NewPostService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("Post Service Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("postServiceError", errorMessage)}, nil

		}

		if err := postService.DeletePostByOwner(req.UserID, postUUID); err != nil {
			errorMessage := fmt.Sprintf("Delete Post Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("deletePostError", errorMessage)}, nil

		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
