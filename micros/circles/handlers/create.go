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
	domain "github.com/red-gold/ts-serverless/micros/circles/dto"
	models "github.com/red-gold/ts-serverless/micros/circles/models"
	service "github.com/red-gold/ts-serverless/micros/circles/services"
)

const followingCircleName = "Following"

// CreateCircleHandle handle create a new circle
func CreateCircleHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.CreateCircleModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse CreateCircleModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	if model.Name == followingCircleName {
		errorMessage := fmt.Sprintf("Can not use 'Following' as a circle name")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("followingCircleNameIsReserved", errorMessage))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[CreateCircleHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	// Create a new circle
	newCircle := &domain.Circle{
		OwnerUserId: currentUser.UserID,
		Name:        model.Name,
		IsSystem:    false,
	}

	// Create service
	circleService, serviceErr := service.NewCircleService(database.Db)
	if serviceErr != nil {
		log.Error("NewCircleService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/circleService", "Error happened while creating circleService!"))
	}

	if err := circleService.SaveCircle(newCircle); err != nil {
		errorMessage := fmt.Sprintf("Save Circle Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveCircle", "Error happened while saving circle!"))
	}

	return c.JSON(fiber.Map{
		"objectId": newCircle.ObjectId.String(),
	})

}

// CreateFollowingHandle handle create a new circle
func CreateFollowingHandle(c *fiber.Ctx) error {

	// params from /circles/following/:userId
	userId := c.Params("userId")
	if userId == "" {
		errorMessage := fmt.Sprintf("User Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("userIdRequired", errorMessage))
	}

	userUUID, uuidErr := uuid.FromString(userId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("userIdIsNotValid", "User id is not valid!"))
	}

	// Create a new circle
	newCircle := &domain.Circle{
		OwnerUserId: userUUID,
		Name:        followingCircleName,
		IsSystem:    true,
	}

	// Create service
	circleService, serviceErr := service.NewCircleService(database.Db)
	if serviceErr != nil {
		log.Error("NewCircleService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/circleService", "Error happened while creating circleService!"))
	}

	if err := circleService.SaveCircle(newCircle); err != nil {
		errorMessage := fmt.Sprintf("Save Circle Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveCircle", "Error happened while saving circle!"))
	}

	return c.JSON(fiber.Map{
		"objectId": newCircle.ObjectId.String(),
	})

}
