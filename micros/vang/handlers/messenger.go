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
	Room  models.RoomModel       `json:"room"`
	Users map[string]interface{} `json:"users"`
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

	// Get the participants profile
	participantsProfile, roomMemberIds, err := getParticipantsProfile(&currentUser, model)
	if err != nil {
		log.Error("[ActivePeerRoom] Error while getting participants profile %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/getParticipantsProfile",
			"Error while getting participants profile"))
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
			DeactiveUsers: []string{roomMemberIds[1]},
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
	}

	roomModel := models.RoomModel{
		ObjectId:      room.ObjectId,
		Members:       room.Members,
		Type:          room.Type,
		ReadDate:      room.ReadDate,
		ReadCount:     room.ReadCount,
		ReadMessageId: room.ReadMessageId,
		DeactiveUsers: room.DeactiveUsers,
		LastMessage:   room.LastMessage,
		MemberCount:   room.MemberCount,
		MessageCount:  room.MessageCount,
		CreatedDate:   room.CreatedDate,
		UpdatedDate:   room.UpdatedDate,
	}

	actionRoomPayload := &SetActiveRoomPayload{
		Room:  roomModel,
		Users: participantsProfile,
	}

	activeRoomAction := Action{
		Type:    "SET_ACTIVE_ROOM",
		Payload: actionRoomPayload,
	}

	if model.ResponseActionType != "" {
		activeRoomAction.Type = model.ResponseActionType
	}
	go dispatchAction(activeRoomAction, currentUser.UserID, getUserInfoReqFromCurrentUser(currentUser))
	return c.JSON(roomModel)
}

// getParticipantsProfile
func getParticipantsProfile(currentUser *types.UserContext, model *models.ActivePeerRoomModel) (map[string]interface{}, []string, error) {

	mappedParticipants := make(map[string]interface{})

	if model.PeerUserId == uuid.Nil && model.SocialName == "" {
		return nil, nil, fmt.Errorf("PeerUserId is nil and SocialName is empty")
	}

	log.Info(fmt.Sprintf("PeerUserId is %s and SocialName is %s", model.PeerUserId, model.SocialName))
	log.Info("CURRENT USER %v", currentUser)

	receptionist := new(models.UserProfileModel)

	var err error
	if model.SocialName != "" {
		receptionist, err = getProfileBySocialName(model.SocialName)
		if err != nil {
			return nil, nil, fmt.Errorf("Get user profile by social name %s", err.Error())
		}

	} else {
		receptionist, err = getUserProfileByID(model.PeerUserId)
		if err != nil {
			return nil, nil, fmt.Errorf("Get user profile by ID %s", err.Error())
		}
	}

	// Map receptionist user
	mappedUser := make(map[string]interface{})
	mappedUser["userId"] = receptionist.ObjectId
	mappedUser["fullName"] = receptionist.FullName
	mappedUser["socialName"] = receptionist.SocialName
	mappedUser["avatar"] = receptionist.Avatar
	mappedUser["banner"] = receptionist.Banner
	mappedUser["tagLine"] = receptionist.TagLine
	mappedUser["lastSeen"] = receptionist.LastSeen
	mappedUser["createdDate"] = receptionist.CreatedDate
	mappedParticipants[receptionist.ObjectId.String()] = mappedUser

	// Map current user
	mappedCurrentUser := make(map[string]interface{})
	mappedCurrentUser["userId"] = currentUser.UserID
	mappedCurrentUser["fullName"] = currentUser.DisplayName
	mappedCurrentUser["socialName"] = currentUser.SocialName
	mappedCurrentUser["avatar"] = currentUser.Avatar
	mappedCurrentUser["tagLine"] = currentUser.TagLine
	mappedCurrentUser["lastSeen"] = utils.UTCNowUnix()
	mappedCurrentUser["createdDate"] = currentUser.CreatedDate
	mappedParticipants[currentUser.UserID.String()] = mappedCurrentUser
	roomMemberIds := []string{currentUser.UserID.String(), receptionist.ObjectId.String()}

	return mappedParticipants, roomMemberIds, nil
}
