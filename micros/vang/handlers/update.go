package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/vang/database"
	domain "github.com/red-gold/ts-serverless/micros/vang/dto"
	models "github.com/red-gold/ts-serverless/micros/vang/models"
	service "github.com/red-gold/ts-serverless/micros/vang/services"
)

// UpdateMessageHandle handle create a new vang
func UpdateMessageHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.MessageModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse MessageModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	messageService, serviceErr := service.NewMessageService(database.Db)
	if serviceErr != nil {
		log.Error("NewMessageService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/messageService", "Error happened while creating messageService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[UpdateMessageHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	updatedMessage := &domain.Message{
		ObjectId:    model.ObjectId,
		OwnerUserId: currentUser.UserID,
		RoomId:      model.RoomId,
		Text:        model.Text,
		CreatedDate: model.CreatedDate,
		UpdatedDate: model.UpdatedDate,
	}

	if err := messageService.UpdateMessageById(updatedMessage); err != nil {
		errorMessage := fmt.Sprintf("Update Message Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updateMessage", "Error happened while updating message!"))
	}
	return c.SendStatus(http.StatusOK)
}

// UpdateMessageHandle handle create a new vang
func UpdateReadMessageHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.UpdateReadMessageModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse UpdateReadMessageModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	roomService, serviceErr := service.NewRoomService(database.Db)
	if serviceErr != nil {
		log.Error("NewRoomService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/roomService", "Error happened while creating roomService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[UpdateReadMessageHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := roomService.UpdateMemberRead(model.RoomId, currentUser.UserID, model.Amount, model.MessageCreatedDate); err != nil {
		errorMessage := fmt.Sprintf("Update Message Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updateMessage", "Error happened while updating message!"))
	}

	return c.SendStatus(http.StatusOK)
}
