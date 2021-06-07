package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	log "github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/vang/database"
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
func ActivePeerRoom(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.ActivePeerRoomModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse SaveMessagesModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[ActivePeerRoom] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}
	roomMemberIds := []string{currentUser.UserID.String(), model.PeerUserId.String()}

	userInfoReq := &UserInfoInReq{
		UserId:      currentUser.UserID,
		Username:    currentUser.Username,
		Avatar:      currentUser.Avatar,
		DisplayName: currentUser.DisplayName,
		SystemRole:  currentUser.SystemRole,
	}

	getProfilesModel := models.GetProfilesModel{
		UserIds: roomMemberIds,
	}

	foundUserProfiles, getPeerUserErr := getProfilesByUserIds(getProfilesModel, userInfoReq)
	if getPeerUserErr != nil {
		errorMessage := fmt.Sprintf("Get profiles %s", getPeerUserErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/getProfiles", "Error happened while reading profiles!"))
	}

	if len(foundUserProfiles) != len(roomMemberIds) {
		errorMessage := fmt.Sprintf("Could not find all profiles (%v)", roomMemberIds)
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/notFoundAllProfile", "Error happened while finding profiles!"))
	}

	// Create service
	roomService, serviceErr := service.NewRoomService(database.Db)
	if serviceErr != nil {
		log.Error("NewRoomService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/roomService", "Error happened while creating roomService!"))
	}
	room, findRoomErr := roomService.FindOneRoomByMembers(roomMemberIds, 0)
	if findRoomErr != nil {
		errorMessage := fmt.Sprintf("Vang find room %s", findRoomErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/findRoom", "Error happened while finding room!"))
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
			log.Error(errorMessage)
			return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveRoom", "Error happened while saving room!"))
		}
		room = &newRoom
	} else {
		// Create service
		messageService, serviceErr := service.NewMessageService(database.Db)
		if serviceErr != nil {
			log.Error("NewMessageService %s", serviceErr.Error())
			return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/messageService", "Error happened while creating messageService!"))
		}
		lastMessages, getMessagesErr := messageService.GetMessageByRoomId(&room.ObjectId, "createdDate", 1, utils.UTCNowUnix(), 0)
		if getMessagesErr != nil {
			errorMessage := fmt.Sprintf("Vang get room messages %s", getMessagesErr.Error())
			log.Error(errorMessage)
			return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/getRoomMessages", "Error happened while reading messages room!"))
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
	return c.JSON(roomModel)
}
