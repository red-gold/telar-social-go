package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/vang/dto"
	models "github.com/red-gold/ts-serverless/micros/vang/models"
	service "github.com/red-gold/ts-serverless/micros/vang/services"
)

type SetActiveRoomPayload struct {
	Room     models.RoomModel       `json:"room"`
	Messages []models.MessageModel  `json:"messages"`
	Users    map[string]interface{} `json:"users"`
}

// ActivePeerRoom handle active peer room
func ActivePeerRoom(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.ActivePeerRoomModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal models.ActivePeerRoomModel %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelUnMarshalError", errorMessage)}, nil
		}

		roomMemberIds := []string{req.UserID.String(), model.PeerUserId.String()}

		userInfoReq := &UserInfoInReq{
			UserId:      req.UserID,
			Username:    req.Username,
			Avatar:      req.Avatar,
			DisplayName: req.DisplayName,
			SystemRole:  req.SystemRole,
		}

		getProfilesModel := models.GetProfilesModel{
			UserIds: roomMemberIds,
		}

		foundUserProfiles, getPeerUserErr := getProfilesByUserIds(getProfilesModel, userInfoReq)
		if getPeerUserErr != nil {
			errorMessage := fmt.Sprintf("Get profiles %s", getPeerUserErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("getProfilesError", errorMessage)}, nil
		}
		if len(foundUserProfiles) != len(roomMemberIds) {
			errorMessage := fmt.Sprintf("Could not find all profiles (%v)", roomMemberIds)
			println(fmt.Sprintf("Could not find all profiles (%v)", roomMemberIds), foundUserProfiles)
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("notFoundAllProfilesError", errorMessage)}, nil
		}

		// Create service
		roomService, serviceErr := service.NewRoomService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("Vang room service %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("roomServiceError", errorMessage)}, nil
		}
		room, findRoomErr := roomService.FindOneRoomByMembers(roomMemberIds, 0)
		if findRoomErr != nil {
			errorMessage := fmt.Sprintf("Vang find room %s", findRoomErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("findRoomError", errorMessage)}, nil
		}

		var roomMessages []models.MessageModel
		if room == nil {
			readDateMap := make(map[string]int64)
			readDateMap[roomMemberIds[0]] = 0
			readDateMap[roomMemberIds[1]] = 0

			readCountMap := make(map[string]int64)
			readCountMap[roomMemberIds[0]] = 0
			readCountMap[roomMemberIds[1]] = 0

			readMessageIdMap := make(map[string]string)
			readMessageIdMap[roomMemberIds[0]] = ""
			readMessageIdMap[roomMemberIds[1]] = ""

			lastMessageMap := make(map[string]interface{})
			newRoom := dto.Room{
				ObjectId:      uuid.Must(uuid.NewV4()),
				Members:       roomMemberIds,
				Type:          0,
				ReadDate:      readDateMap,
				ReadCount:     readCountMap,
				ReadMessageId: readMessageIdMap,
				LastMessage:   lastMessageMap,
				MemberCount:   2,
				MessageCount:  0,
				CreatedDate:   utils.UTCNowUnix(),
				UpdatedDate:   utils.UTCNowUnix(),
			}
			err := roomService.SaveRoom(&newRoom)
			if err != nil {
				errorMessage := fmt.Sprintf("Vang save room %s", err.Error())
				return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("saveRoomError", errorMessage)}, nil
			}
			room = &newRoom
		} else {
			// Create service
			messageService, serviceErr := service.NewMessageService(db)
			if serviceErr != nil {
				errorMessage := fmt.Sprintf("Vang message service %s", serviceErr.Error())
				return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("messageServiceError", errorMessage)}, nil
			}
			lastMessages, getMessagesErr := messageService.GetMessageByRoomId(&room.ObjectId, "createdDate", 1, utils.UTCNowUnix(), 0)
			if getMessagesErr != nil {
				errorMessage := fmt.Sprintf("Vang get room messages %s", getMessagesErr.Error())
				return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("getRoomMessagesError", errorMessage)}, nil
			}
			for _, v := range lastMessages {
				parsedMessage := models.MessageModel{
					ObjectId:    v.ObjectId,
					OwnerUserId: v.OwnerUserId,
					RoomId:      v.RoomId,
					Text:        v.Text,
					CreatedDate: v.CreatedDate,
					UpdatedDate: v.UpdatedDate,
				}
				roomMessages = append(roomMessages, parsedMessage)
			}
		}

		// Map user profiles
		mappedUsers := make(map[string]interface{})
		for _, v := range foundUserProfiles {
			mappedUser := make(map[string]interface{})
			mappedUser["userId"] = v.ObjectId
			mappedUser["fullName"] = v.FullName
			mappedUser["avatar"] = v.Avatar
			mappedUser["banner"] = v.Banner
			mappedUser["tagLine"] = v.TagLine
			mappedUser["lastSeen"] = v.LastSeen
			mappedUser["createdDate"] = v.CreatedDate

			mappedUsers[v.ObjectId.String()] = mappedUser
		}

		roomModel := models.RoomModel{
			ObjectId:      room.ObjectId,
			Members:       room.Members,
			Type:          room.Type,
			ReadDate:      room.ReadDate,
			ReadCount:     room.ReadCount,
			ReadMessageId: room.ReadMessageId,
			LastMessage:   room.LastMessage,
			MemberCount:   room.MemberCount,
			MessageCount:  room.MessageCount,
			CreatedDate:   room.CreatedDate,
			UpdatedDate:   room.UpdatedDate,
		}

		body, marshalError := json.Marshal(&roomModel)
		if marshalError != nil {
			errorMessage := fmt.Sprintf("Marshal room %s", marshalError.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("roomMarshalError", errorMessage)}, nil
		}

		actionRoomPayload := &SetActiveRoomPayload{
			Room:     roomModel,
			Messages: roomMessages,
			Users:    mappedUsers,
		}

		activeRoomAction := Action{
			Type:    "SET_ACTIVE_ROOM",
			Payload: actionRoomPayload,
		}

		go dispatchAction(activeRoomAction, userInfoReq)
		return handler.Response{
			Body:       body,
			StatusCode: http.StatusOK,
		}, nil
	}
}
