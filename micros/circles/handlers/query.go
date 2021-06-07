package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/pkg/parser"
	"github.com/red-gold/telar-core/types"
	utils "github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/circles/database"
	service "github.com/red-gold/ts-serverless/micros/circles/services"
)

type CircleQueryModel struct {
	Search string    `query:"search"`
	Page   int64     `query:"page"`
	Owner  uuid.UUID `query:"owner"`
}

// QueryCircleHandle handle query on circle
func QueryCircleHandle(c *fiber.Ctx) error {

	// Create service
	circleService, serviceErr := service.NewCircleService(database.Db)
	if serviceErr != nil {
		log.Error("NewCircleService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/circleService", "Error happened while creating circleService!"))
	}

	query := new(CircleQueryModel)

	if err := parser.QueryParser(c, query); err != nil {
		log.Error("[QueryCircleHandle] QueryParser %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryParser", "Error happened while parsing query!"))
	}

	circleList, err := circleService.QueryCircle(query.Search, &query.Owner, "created_date", query.Page)
	if err != nil {
		log.Error("Query circle %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryCircle", "Can not query circle!"))
	}

	return c.JSON(circleList)

}

// GetMyCircleHandle handle get authed user circle
func GetMyCircleHandle(c *fiber.Ctx) error {

	// Create service
	circleService, serviceErr := service.NewCircleService(database.Db)
	if serviceErr != nil {
		log.Error("NewCircleService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/circleService", "Error happened while creating circleService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[GetMyCircleHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	circleList, err := circleService.FindByOwnerUserId(currentUser.UserID)
	if err != nil {
		log.Error("[GetMyCircleHandle.circleService.FindByOwnerUserId] %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryParser", "Error happened while finding circle by user id!"))
	}

	return c.JSON(circleList)

}

// GetCircleHandle handle get a circle
func GetCircleHandle(c *fiber.Ctx) error {

	// Create service
	circleService, serviceErr := service.NewCircleService(database.Db)
	if serviceErr != nil {
		log.Error("NewCircleService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/circleService", "Error happened while creating circleService!"))
	}
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

	foundCircle, err := circleService.FindById(circleUUID)
	if err != nil {
		errorMessage := fmt.Sprintf("Find Circle %s - %s", circleUUID.String(), err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("findCircle", "Can not find circle!"))
	}

	return c.JSON(foundCircle)

}
