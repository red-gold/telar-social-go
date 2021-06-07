package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/user-rels/database"
	domain "github.com/red-gold/ts-serverless/micros/user-rels/dto"
	models "github.com/red-gold/ts-serverless/micros/user-rels/models"
	service "github.com/red-gold/ts-serverless/micros/user-rels/services"
)

// UpdateUserRelHandle handle create a new userRel
func UpdateUserRelHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(domain.UserRel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse UserRel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	userRelService, serviceErr := service.NewUserRelService(database.Db)
	if serviceErr != nil {
		log.Error("NewUserRelService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userRelService", "Error happened while creating userRelService!"))
	}

	if err := userRelService.UpdateUserRelById(model); err != nil {
		errorMessage := fmt.Sprintf("Update UserRel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updateUserRelById", "Error happened while updating UserRel!"))
	}
	return c.SendStatus(http.StatusOK)
}

// UpdateRelCirclesHandle handle create a new userRel
func UpdateRelCirclesHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.RelCirclesModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse RelCirclesModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	userRelService, serviceErr := service.NewUserRelService(database.Db)
	if serviceErr != nil {
		log.Error("NewUserRelService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userRelService", "Error happened while creating userRelService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[UpdateRelCirclesHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := userRelService.UpdateRelCircles(currentUser.UserID, model.RightId, model.CircleIds); err != nil {
		errorMessage := fmt.Sprintf("Update UserRel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updateRelCircles", "Error happened while updating UserRel!"))
	}
	return c.SendStatus(http.StatusOK)
}
