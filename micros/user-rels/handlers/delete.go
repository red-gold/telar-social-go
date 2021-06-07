package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/user-rels/database"
	service "github.com/red-gold/ts-serverless/micros/user-rels/services"
)

// DeleteUserRelHandle handle delete a userRel
func DeleteUserRelHandle(c *fiber.Ctx) error {

	// params from /user-rels/:userRelId
	userRelId := c.Params("userRelId")
	if userRelId == "" {
		errorMessage := fmt.Sprintf("UserRel Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("userRelIdRequired", errorMessage))

	}

	userRelUUID, uuidErr := uuid.FromString(userRelId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("userRelIdIsNotValid", "user rel id is not valid!"))
	}

	// Create service
	userRelService, serviceErr := service.NewUserRelService(database.Db)
	if serviceErr != nil {
		log.Error("NewUserRelService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userRelService", "Error happened while creating userRelService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[DeleteUserRelHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := userRelService.DeleteUserRelByOwner(currentUser.UserID, userRelUUID); err != nil {
		errorMessage := fmt.Sprintf("Delete UserRel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/deleteUserRel", "Error happened while removing user-rel!"))
	}

	return c.SendStatus(http.StatusOK)
}

// UnfollowHandle handle delete a userRel
func UnfollowHandle(c *fiber.Ctx) error {

	// params from /user-rels/unfollow/:userId
	userFollowingId := c.Params("userId")
	if userFollowingId == "" {
		errorMessage := fmt.Sprintf("User Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("userIdRequired", errorMessage))
	}

	userFollowingUUID, uuidErr := uuid.FromString(userFollowingId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("userFollowingUUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("userFollowingIdIsNotValid", "user following id is not valid!"))
	}

	// Create service
	userRelService, serviceErr := service.NewUserRelService(database.Db)
	if serviceErr != nil {
		log.Error("NewUserRelService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userRelService", "Error happened while creating userRelService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[UnfollowHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := userRelService.UnfollowUser(currentUser.UserID, userFollowingUUID); err != nil {
		errorMessage := fmt.Sprintf("Delete UserRel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/unfollowUser", "Error happened while removing user-rel!"))
	}

	// Decrease user follow count
	go increaseUserFollowCount(currentUser.UserID, -1, getUserInfoReq(c))
	// Decrease user follower count
	go increaseUserFollowerCount(userFollowingUUID, -1, getUserInfoReq(c))

	return c.SendStatus(http.StatusOK)
}

// DeleteCircle handle delete a userRel
func DeleteCircle(c *fiber.Ctx) error {

	// params from /user-rels/circle/:circleId
	circleId := c.Params("circleId")
	if circleId == "" {
		errorMessage := fmt.Sprintf("Circle Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("circleIdRequired", errorMessage))
	}

	// Create service
	userRelService, serviceErr := service.NewUserRelService(database.Db)
	if serviceErr != nil {
		log.Error("NewUserRelService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userRelService", "Error happened while creating userRelService!"))
	}

	if err := userRelService.DeleteCircle(circleId); err != nil {
		errorMessage := fmt.Sprintf("Delete circle from user-rel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/deleteCircle", "Error happened while removing circle!"))
	}
	return c.SendStatus(http.StatusOK)
}
