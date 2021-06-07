package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/circles/database"
	domain "github.com/red-gold/ts-serverless/micros/circles/dto"
	service "github.com/red-gold/ts-serverless/micros/circles/services"
)

// UpdateCircleHandle handle create a new circle
func UpdateCircleHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(domain.Circle)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse Circle Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	if model.Name == "" {
		errorMessage := fmt.Sprintf("Circle name can not be empty.")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("circleNameIsRequired", errorMessage))
	}

	// Create service
	circleService, serviceErr := service.NewCircleService(database.Db)
	if serviceErr != nil {
		log.Error("NewCircleService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/circleService", "Error happened while creating circleService!"))
	}

	if err := circleService.UpdateCircleById(model); err != nil {
		errorMessage := fmt.Sprintf("Update Circle Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("updateCircle", "Can not update circle!"))
	}

	return c.SendStatus(http.StatusOK)

}
