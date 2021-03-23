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
	Room     models.RoomModel         `json:"room"`
	Members  []models.RoomMemberModel `json:"members"`
	Messages []models.MessageModel    `json:"messages"`
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

		peerUserProfile, getPeerUserErr := getUserProfileByID(model.PeerUserId)
		if getPeerUserErr != nil {
			errorMessage := fmt.Sprintf("Get peer user profile %s", getPeerUserErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("getPeerUserProfileError", errorMessage)}, nil
		}
		if peerUserProfile == nil {
			errorMessage := fmt.Sprintf("Peer user does not exist (%s)", model.PeerUserId)
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("peerUseNotExistError", errorMessage)}, nil
		}

		// Create service
		roomService, serviceErr := service.NewRoomService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("Vang room service %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("roomServiceError", errorMessage)}, nil
		}
		roomMembers := []string{req.UserID.String(), model.PeerUserId.String()}
		room, findRoomErr := roomService.FindOneRoomByMembers(roomMembers, 0)
		if findRoomErr != nil {
			errorMessage := fmt.Sprintf("Vang find room %s", findRoomErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("findRoomError", errorMessage)}, nil
		}

		var roomMessages []models.MessageModel
		if room == nil {
			seenMap := make(map[string]int64)
			seenMap[roomMembers[0]] = utils.UTCNowUnix()
			seenMap[roomMembers[1]] = utils.UTCNowUnix()
			newRoom := dto.Room{
				ObjectId:    uuid.Must(uuid.NewV4()),
				Members:     roomMembers,
				Type:        0,
				Seen:        seenMap,
				CreatedDate: utils.UTCNowUnix(),
				UpdatedDate: utils.UTCNowUnix(),
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
			lastMessages, getMessagesErr := messageService.GetMessageByRoomId(&room.ObjectId, "createdDate", 1, utils.UTCNowUnix())
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

		roomModel := models.RoomModel{
			ObjectId:    room.ObjectId,
			Members:     room.Members,
			Type:        room.Type,
			Seen:        room.Seen,
			CreatedDate: room.CreatedDate,
			UpdatedDate: room.UpdatedDate,
		}
		body, marshalError := json.Marshal(&roomModel)
		if marshalError != nil {
			errorMessage := fmt.Sprintf("Marshal room %s", marshalError.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("roomMarshalError", errorMessage)}, nil
		}

		peerUser := models.RoomMemberModel{
			ObjectId: peerUserProfile.ObjectId,
			FullName: peerUserProfile.FullName,
			Avatar:   peerUserProfile.Avatar,
		}

		reqUser := models.RoomMemberModel{
			ObjectId: req.UserID,
			FullName: req.DisplayName,
			Avatar:   req.Avatar,
		}

		actionRoomPayload := &SetActiveRoomPayload{
			Room: roomModel,
			Members: []models.RoomMemberModel{
				reqUser,
				peerUser,
			},
			Messages: roomMessages,
		}

		activeRoomAction := Action{
			Type:    "SET_ACTIVE_ROOM",
			Payload: actionRoomPayload,
		}

		userInfoReq := &UserInfoInReq{
			UserId:      req.UserID,
			Username:    req.Username,
			Avatar:      req.Avatar,
			DisplayName: req.DisplayName,
			SystemRole:  req.SystemRole,
		}
		go dispatchAction(activeRoomAction, userInfoReq)
		return handler.Response{
			Body:       body,
			StatusCode: http.StatusOK,
		}, nil
	}
}
