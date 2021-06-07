package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	log "github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/vang/database"
	service "github.com/red-gold/ts-serverless/micros/vang/services"
)

// DeleteMessageHandle handle delete a Message
func DeleteMessageHandle(c *fiber.Ctx) error {

	// params from /message/id/:messageId
	messageId := c.Params("messageId")
	if messageId == "" {
		errorMessage := fmt.Sprintf("Message Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("messageIdRequired", errorMessage))
	}

	messageUUID, uuidErr := uuid.FromString(messageId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("messageIdIsNotValid", "Message id is not valid!"))
	}

	// Create service
	messageService, serviceErr := service.NewMessageService(database.Db)
	if serviceErr != nil {
		log.Error("NewMessageService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/messageService", "Error happened while creating messageService!"))

	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[DeleteMessageHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := messageService.DeleteMessageByOwner(currentUser.UserID, messageUUID); err != nil {
		errorMessage := fmt.Sprintf("Delete Message Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/deleteMessage", "Error happened while removing message!"))
	}

	return c.SendStatus(http.StatusOK)
}
