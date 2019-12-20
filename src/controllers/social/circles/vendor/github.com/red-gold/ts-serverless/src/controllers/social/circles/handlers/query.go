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
	service "github.com/red-gold/ts-serverless/src/services/social"
	uuid "github.com/satori/go.uuid"
)

// QueryCircleHandle handle query on circle
func QueryCircleHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		circleService, serviceErr := service.NewCircleService(db)
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
		circleList, err := circleService.QueryCircle(searchParam, ownerUserId, "created_date", int64(page))
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(circleList)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling circleList: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("circleListMarshalError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// GetMyCircleHandle handle get authed user circle
func GetMyCircleHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		circleService, serviceErr := service.NewCircleService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		circleList, err := circleService.FindByOwnerUserId(req.UserID)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(circleList)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling circleList: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("circleListMarshalError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// GetCircleHandle handle get a circle
func GetCircleHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		circleService, serviceErr := service.NewCircleService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}
		circleId := req.GetParamByName("circleId")
		circleUUID, uuidErr := uuid.FromString(circleId)
		if uuidErr != nil {
			return handler.Response{StatusCode: http.StatusBadRequest,
					Body: utils.MarshalError("parseUUIDError", "Can not parse circle id!")},
				nil
		}

		foundCircle, err := circleService.FindById(circleUUID)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(foundCircle)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling circleModel: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: []byte(utils.MarshalError("marshalCircleModel", errorMessage))},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}
