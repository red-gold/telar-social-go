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
	service "github.com/red-gold/ts-serverless/micros/user-rels/services"
	uuid "github.com/satori/go.uuid"
)

// QueryUserRelHandle handle query on userRel
func QueryUserRelHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		userRelService, serviceErr := service.NewUserRelService(db)
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

		var ownerUserId *uuid.UUID = nil
		if ownerUserIdParam != "" {

			parsedUUID, uuidErr := uuid.FromString(ownerUserIdParam)

			if uuidErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, uuidErr
			}

			ownerUserId = &parsedUUID
		}

		page := 0
		if pageParam != "" {
			var strErr error
			page, strErr = strconv.Atoi(pageParam)
			if strErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, strErr
			}
		}
		userRelList, err := userRelService.QueryUserRel(searchParam, ownerUserId, "created_date", int64(page))
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(userRelList)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling userRelList: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("userRelListMarshalError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// GetUserRelHandle handle get a userRel
func GetUserRelHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		userRelService, serviceErr := service.NewUserRelService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}
		userRelId := req.GetParamByName("userRelId")
		userRelUUID, uuidErr := uuid.FromString(userRelId)
		if uuidErr != nil {
			return handler.Response{StatusCode: http.StatusBadRequest,
					Body: utils.MarshalError("parseUUIDError", "Can not parse userRel id!")},
				nil
		}

		foundUserRel, err := userRelService.FindById(userRelUUID)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(foundUserRel)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling userRelModel: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: []byte(utils.MarshalError("marshalUserRelModel", errorMessage))},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// GetFollowersHandle handle get auth user followers
func GetFollowersHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		userRelService, serviceErr := service.NewUserRelService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		followers, err := userRelService.GetFollowers(req.UserID)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(followers)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling userRelModel: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: []byte(utils.MarshalError("marshalUserRelModel", errorMessage))},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// GetFollowingHandle handle get auth user following
func GetFollowingHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		userRelService, serviceErr := service.NewUserRelService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		followers, err := userRelService.GetFollowing(req.UserID)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(followers)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling userRelModel: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: []byte(utils.MarshalError("marshalUserRelModel", errorMessage))},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}
