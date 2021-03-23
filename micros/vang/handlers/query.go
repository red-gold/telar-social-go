package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	uuid "github.com/gofrs/uuid"
	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	utils "github.com/red-gold/telar-core/utils"
	service "github.com/red-gold/ts-serverless/micros/vang/services"
)

// QueryMessagesHandle handle query on vang
func QueryMessagesHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		vangService, serviceErr := service.NewMessageService(db)
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
		roomIdParam := query.Get("roomId")
		pageParam := query.Get("page")
		lteParam := query.Get("lte")

		var roomId *uuid.UUID = nil
		if roomIdParam != "" {

			parsedUUID, uuidErr := uuid.FromString(roomIdParam)

			if uuidErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, uuidErr
			}

			roomId = &parsedUUID
		} else {
			errorMessage := fmt.Sprintf("Room id can not be empty.")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("roomIdEmptyError", errorMessage)},
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

		lte := 0
		if lteParam != "" {
			var strErr error
			lte, strErr = strconv.Atoi(lteParam)
			if strErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, strErr
			}
		}

		vangList, err := vangService.GetMessageByRoomId(roomId, "createdDate", int64(page), int64(lte))
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(vangList)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling vangList: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("vangListMarshalError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}
