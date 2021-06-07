package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/circles/database"
	service "github.com/red-gold/ts-serverless/micros/circles/services"
)

// DeleteCircleHandle handle delete a circle
func DeleteCircleHandle(c *fiber.Ctx) error {

	// params from /circles/:circleId
	circleId := c.Params("circleId")
	if circleId == "" {
		errorMessage := fmt.Sprintf("Circle Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("circleIdRequired", errorMessage))
	}

	circleUUID, uuidErr := uuid.FromString(circleId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("circleIdIsNotValid", "Circle id is not valid!"))

	}

	// Create service
	circleService, serviceErr := service.NewCircleService(database.Db)
	if serviceErr != nil {
		log.Error("NewCircleService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/circleService", "Error happened while creating circleService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[DeleteCircleHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := circleService.DeleteCircleByOwner(currentUser.UserID, circleUUID); err != nil {
		errorMessage := fmt.Sprintf("Delete Circle Error %s - %s", circleUUID.String(), err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("deleteCircle", "Can not delete circle!"))
	}

	return c.SendStatus(http.StatusOK)
}
