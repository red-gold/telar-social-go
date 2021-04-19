package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	uuid "github.com/gofrs/uuid"
	handler "github.com/openfaas-incubator/go-function-sdk"
	"github.com/red-gold/telar-core/pkg/log"
	server "github.com/red-gold/telar-core/server"
	utils "github.com/red-gold/telar-core/utils"
	models "github.com/red-gold/ts-serverless/micros/vang/models"
	service "github.com/red-gold/ts-serverless/micros/vang/services"
)

// QueryMessagesHandle handle query on vang
func QueryMessagesHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Parse model object
		var model models.QueryMessageModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal models.QueryMessageModel %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelUnmarshalError", errorMessage)}, nil
		}

		if model.ReqUserId != req.UserID {
			errorMessage := fmt.Sprintf("Request user id is not equal.")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("reqUserIdNotEqualError", errorMessage)},
				nil
		}

		// Create service
		vangService, serviceErr := service.NewMessageService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		if model.RoomId == uuid.Nil {
			errorMessage := fmt.Sprintf("Room id can not be empty.")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("roomIdEmptyError", errorMessage)},
				nil
		}

		vangList, err := vangService.GetMessageByRoomId(&model.RoomId, "createdDate", model.Page, model.Lte, model.Gte)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(vangList)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling message list: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("messageListMarshalError", errorMessage)},
				marshalErr

		}

		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// GetUserRooms handle active peer room
func GetUserRooms(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Parse model object
		var model models.GetUserRoomsModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal models.ActivePeerRoomModel %s", err.Error())
			log.Error(errorMessage)
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelUnMarshalError", errorMessage)}, nil
		}

		// Create service
		roomService, serviceErr := service.NewRoomService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("Vang room service %s", serviceErr.Error())
			log.Error(errorMessage)
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("roomServiceError", errorMessage)}, nil
		}

		rooms, findRoomErr := roomService.GetRoomsByUserId(model.UserId.String(), 0)
		if findRoomErr != nil {
			errorMessage := fmt.Sprintf("Vang find room %s", findRoomErr.Error())
			log.Error(errorMessage)
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("findRoomError", errorMessage)}, nil
		}

		if len(rooms) == 0 {
			return handler.Response{
				Body:       []byte("{\n  \"rooms\": {},\n  \"roomIds\": []\n}"),
				StatusCode: http.StatusOK,
			}, nil
		}

		// Use map to record duplicates as we find them.
		encountered := map[string]bool{}

		var allMembers []string

		// Map to response model
		var resRooms models.ResUserRoomModel
		resRooms.Rooms = make(map[string]interface{})
		for _, v := range rooms {
			roomId := v.ObjectId.String()
			mappedRoom := make(map[string]interface{})
			mappedRoom["objectId"] = roomId
			mappedRoom["members"] = v.Members
			mappedRoom["type"] = v.Type
			mappedRoom["readDate"] = v.ReadDate
			mappedRoom["readCount"] = v.ReadCount
			mappedRoom["ReadMessageId"] = v.ReadMessageId
			mappedRoom["lastMessage"] = v.LastMessage
			mappedRoom["memberCount"] = v.MemberCount
			mappedRoom["messageCount"] = v.MessageCount
			mappedRoom["createdDate"] = v.CreatedDate
			mappedRoom["updatedDate"] = v.UpdatedDate

			resRooms.Rooms[roomId] = mappedRoom
			resRooms.RoomIds = append(resRooms.RoomIds, roomId)

			// Merge members into a single array
			for _, v := range v.Members[:2] {
				if encountered[v] != true {
					encountered[v] = true
					allMembers = append(allMembers, v)
				}
			}
		}

		dispatchProfileModel := models.DispatchProfilesModel{
			UserIds:   allMembers,
			ReqUserId: model.UserId,
		}

		userInfoInReq := &UserInfoInReq{
			UserId:      req.UserID,
			Username:    req.Username,
			Avatar:      req.Avatar,
			DisplayName: req.DisplayName,
			SystemRole:  req.SystemRole,
		}
		go dispatchProfileByUserIds(dispatchProfileModel, userInfoInReq)

		body, marshalError := json.Marshal(&resRooms)
		if marshalError != nil {
			errorMessage := fmt.Sprintf("Marshal rooms %s", marshalError.Error())
			log.Error(errorMessage)
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("roomsMarshalError", errorMessage)}, nil
		}

		return handler.Response{
			Body:       body,
			StatusCode: http.StatusOK,
		}, nil
	}
}
