package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	utils "github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/vang/database"
	models "github.com/red-gold/ts-serverless/micros/vang/models"
	service "github.com/red-gold/ts-serverless/micros/vang/services"
)

// QueryMessagesHandle handle query on vang
func QueryMessagesHandle(c *fiber.Ctx) error {

	// Parse model object
	model := new(models.QueryMessageModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse SaveMessagesModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[QueryMessagesHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if model.ReqUserId != currentUser.UserID {
		errorMessage := fmt.Sprintf("Request user id is not equal.")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("reqUserIdNotEqual", errorMessage))
	}

	// Create service
	vangService, serviceErr := service.NewMessageService(database.Db)
	if serviceErr != nil {
		log.Error("NewMessageService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/messageService", "Error happened while creating messageService!"))
	}

	if model.RoomId == uuid.Nil {
		errorMessage := fmt.Sprintf("Room id can not be empty.")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("roomIdIsRequired", errorMessage))

	}

	vangList, err := vangService.GetMessageByRoomId(&model.RoomId, "createdDate", model.Page, model.Lte, model.Gte)
	if err != nil {
		log.Error("[QueryMessagesHandle.vangService.GetMessageByRoomId] %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/getMessages", "Error happened while reading messages!"))
	}

	return c.JSON(vangList)
}

// GetUserRooms handle active peer room
func GetUserRooms(c *fiber.Ctx) error {

	// Parse model object
	model := new(models.GetUserRoomsModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse SaveMessagesModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	roomService, serviceErr := service.NewRoomService(database.Db)
	if serviceErr != nil {
		log.Error("NewRoomService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/roomService", "Error happened while creating roomService!"))
	}

	rooms, findRoomErr := roomService.GetRoomsByUserId(model.UserId.String(), 0)
	if findRoomErr != nil {
		log.Error("[GetUserRooms.roomService.GetRoomsByUserId] %s", findRoomErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/findRoom", "Error happened while finding room!"))
	}

	if len(rooms) == 0 {
		c.JSON(fiber.Map{
			"rooms":   fiber.Map{},
			"roomIds": []string{},
		})
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

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[GetUserRooms] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	userInfoInReq := &UserInfoInReq{
		UserId:      currentUser.UserID,
		Username:    currentUser.Username,
		Avatar:      currentUser.Avatar,
		DisplayName: currentUser.DisplayName,
		SystemRole:  currentUser.SystemRole,
	}
	go dispatchProfileByUserIds(dispatchProfileModel, userInfoInReq)

	return c.JSON(resRooms)
}
