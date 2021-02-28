package handlers

import (
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	service "github.com/red-gold/ts-serverless/micros/votes/services"
	uuid "github.com/satori/go.uuid"
)

// DeleteVoteHandle handle delete a Vote
func DeleteVoteHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /votes/id/:voteId
		voteId := req.GetParamByName("voteId")
		if voteId == "" {
			errorMessage := fmt.Sprintf("Vote Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("voteIdRequired", errorMessage)}, nil
		}
		fmt.Printf("\n Vote ID: %s", voteId)
		voteUUID, uuidErr := uuid.FromString(voteId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}
		fmt.Printf("\n Vote UUID: %s", voteUUID)
		// Create service
		voteService, serviceErr := service.NewVoteService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("Vote Service Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("voteServiceError", errorMessage)}, nil

		}

		if err := voteService.DeleteVoteByOwner(req.UserID, voteUUID); err != nil {
			errorMessage := fmt.Sprintf("Delete Vote Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("deleteVoteError", errorMessage)}, nil

		}

		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// DeleteVoteByPostIdHandle handle delete a Vote but postId
func DeleteVoteByPostIdHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /Votes/post/:postId
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
		voteService, serviceErr := service.NewVoteService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("Vote Service Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("VoteServiceError", errorMessage)}, nil

		}

		if err := voteService.DeleteVotesByPostId(req.UserID, PostUUID); err != nil {
			errorMessage := fmt.Sprintf("Delete Vote Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("deleteVoteError", errorMessage)}, nil
		}

		// Create user headers for http request
		userHeaders := make(map[string][]string)
		userHeaders["uid"] = []string{req.UserID.String()}
		userHeaders["email"] = []string{req.Username}
		userHeaders["avatar"] = []string{req.Avatar}
		userHeaders["displayName"] = []string{req.DisplayName}
		userHeaders["role"] = []string{req.SystemRole}

		postIndexURL := fmt.Sprintf("/posts/score/-1/%s", postId)
		_, postIndexErr := functionCall(http.MethodPut, []byte(""), postIndexURL, userHeaders)

		if postIndexErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError,
					Body: utils.MarshalError("decreasePostScoreError",
						fmt.Sprintf("Cannot decrease post score! error: %s", postIndexErr.Error()))},
				nil
		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
