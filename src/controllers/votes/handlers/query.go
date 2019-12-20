package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	utils "github.com/red-gold/telar-core/utils"
	models "github.com/red-gold/ts-serverless/src/models/votes"
	service "github.com/red-gold/ts-serverless/src/services/votes"
	uuid "github.com/satori/go.uuid"
)

// GetVotesByPostIdHandle handle query on vote
func GetVotesByPostIdHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		voteService, serviceErr := service.NewVoteService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		var query *url.Values
		if len(req.QueryString) > 0 {
			q, err := url.ParseQuery(string(req.QueryString))
			if err != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, err
			}
			query = &q
		}
		postIdParam := query.Get("postId")
		pageParam := query.Get("page")

		var postId *uuid.UUID = nil
		if postIdParam != "" {

			parsedUUID, uuidErr := uuid.FromString(postIdParam)

			if uuidErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, uuidErr
			}

			postId = &parsedUUID
		} else {
			errorMessage := fmt.Sprintf("Post id can not be empty.")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("postIdEmptyError", errorMessage)},
				nil
		}

		page := 0
		if pageParam != "" {
			var strErr error
			page, strErr = strconv.Atoi(pageParam)
			if strErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, strErr
			}
		}

		voteList, err := voteService.GetVoteByPostId(postId, "created_date", int64(page))
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(voteList)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling voteList: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("voteListMarshalError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// GetVoteHandle handle get a vote
func GetVoteHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		voteService, serviceErr := service.NewVoteService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}
		voteId := req.GetParamByName("voteId")
		voteUUID, uuidErr := uuid.FromString(voteId)
		if uuidErr != nil {
			return handler.Response{StatusCode: http.StatusBadRequest,
					Body: utils.MarshalError("parseUUIDError", "Can not parse vote id!")},
				nil
		}

		foundVote, err := voteService.FindById(voteUUID)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		voteModel := models.VoteModel{
			ObjectId:         foundVote.ObjectId,
			OwnerUserId:      foundVote.OwnerUserId,
			PostId:           foundVote.PostId,
			OwnerDisplayName: foundVote.OwnerDisplayName,
			OwnerAvatar:      foundVote.OwnerAvatar,
			CreatedDate:      foundVote.CreatedDate,
			TypeId:           foundVote.TypeId,
		}

		body, marshalErr := json.Marshal(voteModel)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("{error: 'Error while marshaling voteModel: %s'}",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: []byte(errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}
