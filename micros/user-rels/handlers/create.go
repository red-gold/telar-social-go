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
	socialModels "github.com/red-gold/ts-serverless/micros/user-rels/models"
	service "github.com/red-gold/ts-serverless/micros/user-rels/services"
)

// CreateUserRelHandle handle create a new userRel
func CreateUserRelHandle(c *fiber.Ctx) error {

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

	if err := userRelService.SaveUserRel(model); err != nil {
		errorMessage := fmt.Sprintf("Save UserRel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveUserRel", "Error happened while saving UserRel!"))
	}

	return c.JSON(fiber.Map{
		"objectId": model.ObjectId.String(),
	})

}

//FollowHandle handle create a new userRel
func FollowHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(socialModels.FollowModel)
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

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[FollowHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	// Left User Meta
	leftUserMeta := domain.UserRelMeta{
		UserId:   currentUser.UserID,
		FullName: currentUser.DisplayName,
		Avatar:   currentUser.Avatar,
	}

	// Right User Meta
	rightUserMeta := domain.UserRelMeta{
		UserId:   model.RightUser.UserId,
		FullName: model.RightUser.FullName,
		Avatar:   model.RightUser.Avatar,
	}

	if err := userRelService.FollowUser(leftUserMeta, rightUserMeta, model.CircleIds, []string{"status:follow"}); err != nil {
		errorMessage := fmt.Sprintf("Save UserRel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveUserRel", "Error happened while saving UserRel!"))
	}

	// Create notification
	go sendFollowNotification(model, getUserInfoReq(c))
	// Increase user follow count
	go increaseUserFollowCount(currentUser.UserID, 1, getUserInfoReq(c))
	// Increase user follower count
	go increaseUserFollowerCount(model.RightUser.UserId, 1, getUserInfoReq(c))

	return c.SendStatus(http.StatusOK)
}
