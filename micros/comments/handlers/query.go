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
	models "github.com/red-gold/ts-serverless/micros/comments/models"
	service "github.com/red-gold/ts-serverless/micros/comments/services"
	uuid "github.com/satori/go.uuid"
)

// QueryCommentHandle handle query on comment
func QueryCommentHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		commentService, serviceErr := service.NewCommentService(db)
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
		searchParam := query.Get("search")
		pageParam := query.Get("page")
		ownerUserIdParam := query.Get("owner")
		commentTypeIdParam := query.Get("type")

		var ownerUserId *uuid.UUID = nil
		if ownerUserIdParam != "" {

			parsedUUID, uuidErr := uuid.FromString(ownerUserIdParam)

			if uuidErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, uuidErr
			}

			ownerUserId = &parsedUUID
		}

		var commentTypeId *int = nil
		if commentTypeIdParam != "" {

			parsedType, strErr := strconv.Atoi(commentTypeIdParam)
			if strErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, strErr
			}
			commentTypeId = &parsedType
		}
		page := 0
		if pageParam != "" {
			var strErr error
			page, strErr = strconv.Atoi(pageParam)
			if strErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, strErr
			}
		}
		commentList, err := commentService.QueryComment(searchParam, ownerUserId, commentTypeId, "created_date", int64(page))
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(commentList)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling commentList: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("commentListMarshalError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// GetCommentsByPostIdHandle handle query on comment
func GetCommentsByPostIdHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		commentService, serviceErr := service.NewCommentService(db)
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

		commentList, err := commentService.GetCommentByPostId(postId, "created_date", int64(page))
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(commentList)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling commentList: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("commentListMarshalError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// GetCommentHandle handle get a comment
func GetCommentHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		commentService, serviceErr := service.NewCommentService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}
		commentId := req.GetParamByName("commentId")
		commentUUID, uuidErr := uuid.FromString(commentId)
		if uuidErr != nil {
			return handler.Response{StatusCode: http.StatusBadRequest,
					Body: utils.MarshalError("parseUUIDError", "Can not parse comment id!")},
				nil
		}

		foundComment, err := commentService.FindById(commentUUID)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		// No comment found
		if foundComment == nil {
			return handler.Response{
				Body:       []byte(nil),
				StatusCode: http.StatusOK,
			}, nil
		}
		commentModel := models.CommentModel{
			ObjectId:         foundComment.ObjectId,
			OwnerUserId:      foundComment.OwnerUserId,
			PostId:           foundComment.PostId,
			Score:            foundComment.Score,
			Text:             foundComment.Text,
			OwnerDisplayName: foundComment.OwnerDisplayName,
			OwnerAvatar:      foundComment.OwnerAvatar,
			Deleted:          foundComment.Deleted,
			DeletedDate:      foundComment.DeletedDate,
			CreatedDate:      foundComment.CreatedDate,
			LastUpdated:      foundComment.LastUpdated,
		}

		body, marshalErr := json.Marshal(commentModel)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("{error: 'Error while marshaling commentModel: %s'}",
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
