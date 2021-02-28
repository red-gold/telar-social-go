package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	domain "github.com/red-gold/ts-serverless/micros/votes/dto"
	models "github.com/red-gold/ts-serverless/micros/votes/models"
	service "github.com/red-gold/ts-serverless/micros/votes/services"
)

// UpdateVoteHandle handle create a new vote
func UpdateVoteHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.VoteModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		// Create service
		voteService, serviceErr := service.NewVoteService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		updatedVote := &domain.Vote{
			ObjectId:         model.ObjectId,
			PostId:           model.PostId,
			OwnerDisplayName: req.DisplayName,
			OwnerAvatar:      req.Avatar,
			CreatedDate:      model.CreatedDate,
		}

		if err := voteService.UpdateVoteById(updatedVote); err != nil {
			errorMessage := fmt.Sprintf("Update Vote Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("updateVoteError", errorMessage)}, nil
		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
