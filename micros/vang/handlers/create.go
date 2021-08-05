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

// CreateMessageHandle handle create a new vang
func SaveMessages(c *fiber.Ctx) error {

	log.Info("[SaveMessages] hit ...")
	// Create the model object
	model := new(models.SaveMessagesModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse SaveMessagesModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Room service
	roomService, roomServiceErr := service.NewRoomService(database.Db)
	if roomServiceErr != nil {
		log.Error("NewRoomService %s", roomServiceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/roomService", "Error happened while creating roomService!"))
	}

	log.Info("[SaveMessages] saving message model.RoomId: %s", model.RoomId)

	// Message service
	messageService, messageServiceErr := service.NewMessageService(database.Db)
	if messageServiceErr != nil {
		log.Error("NewMessageService %s", messageServiceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/messageService", "Error happened while creating messageService!"))
	}

	log.Info("[SaveMessages] message saves")

	// Map message model to DTO
	var messages []dto.Message
	for _, v := range model.Messages {
		newMessage := dto.Message{
			ObjectId:    v.ObjectId,
			OwnerUserId: v.OwnerUserId,
			RoomId:      v.RoomId,
			Text:        v.Text,
			CreatedDate: utils.UTCNowUnix(),
			UpdatedDate: utils.UTCNowUnix(),
		}
		messages = append(messages, newMessage)
	}

	var maxDate int64
	var lastMessage *dto.Message
	// Get last messsage
	for _, v := range messages {
		if v.CreatedDate > maxDate {
			maxDate = v.CreatedDate
			lastMessage = &v
		}
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[SaveMessages] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	log.Info("[SaveMessages] currentUser: %s", currentUser.UserID)

	log.Info("[SaveMessages] check deactive peer id %s", model.DeactivePeerId)
	if model.DeactivePeerId != uuid.Nil && model.DeactivePeerId != currentUser.UserID {
		// Active peer id
		roomErr := roomService.ActiveAllPeerRoom(lastMessage.RoomId, []string{model.DeactivePeerId.String(), currentUser.UserID.String()}, model.DeactivePeerId)
		if roomErr != nil {
			log.Error("[SaveMessages] GetPeerRoom %s", roomErr.Error())
			return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/roomService", "Error happened while getting room!"))
		}
	}
	log.Info("[SaveMessages] check deactive peer id end")

	// Increase room message count
	go func(currentUser types.UserContext) {
		log.Info("[SaveMessages] updating message meta")
		err := roomService.UpdateMessageMeta(model.RoomId, int64(len(messages)), lastMessage.CreatedDate, lastMessage.Text, currentUser.UserID.String())
		if err != nil {
			errorMessage := fmt.Sprintf("vang IncreaseMessageCount %s", err.Error())
			println(errorMessage)
		}
		log.Info("[SaveMessages] updating message meta end")
	}(currentUser)
	log.Info("[SaveMessages] Saving message")

	if err := messageService.SaveManyMessages(messages); err != nil {
		errorMessage := fmt.Sprintf("Save many messages %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveMessage", "Error happened while saving message!"))
	}

	return c.SendStatus(http.StatusOK)
}
