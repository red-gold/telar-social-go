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
	"github.com/red-gold/ts-serverless/micros/user-rels/database"
	service "github.com/red-gold/ts-serverless/micros/user-rels/services"
)

type UserRelQueryModel struct {
	Search string    `query:"search"`
	Page   int64     `query:"page"`
	Owner  uuid.UUID `query:"owner"`
}

// QueryUserRelHandle handle query on userRel
func QueryUserRelHandle(c *fiber.Ctx) error {

	// Create service
	userRelService, serviceErr := service.NewUserRelService(database.Db)
	if serviceErr != nil {
		log.Error("NewUserRelService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userRelService", "Error happened while creating userRelService!"))
	}

	query := new(UserRelQueryModel)

	if err := parser.QueryParser(c, query); err != nil {
		log.Error("[QueryUserRelHandle] QueryParser %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryParser", "Error happened while parsing query!"))
	}

	userRelList, err := userRelService.QueryUserRel(query.Search, &query.Owner, "created_date", query.Page)
	if err != nil {
		log.Error("[QueryUserRelHandle.userRelService.QueryUserRel] %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/queryUserRel", "Error happened while reading followers!"))
	}

	return c.JSON(userRelList)
}

// GetUserRelHandle handle get a userRel
func GetUserRelHandle(c *fiber.Ctx) error {

	// Create service
	userRelService, serviceErr := service.NewUserRelService(database.Db)
	if serviceErr != nil {
		log.Error("NewUserRelService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userRelService", "Error happened while creating userRelService!"))
	}
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

	foundUserRel, err := userRelService.FindById(userRelUUID)
	if err != nil {
		log.Error("[GetUserRelHandle.userRelService.FindById] %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/findById", "Error happened while reading followers!"))
	}

	return c.JSON(foundUserRel)
}

// GetFollowersHandle handle get auth user followers
func GetFollowersHandle(c *fiber.Ctx) error {

	// Create service
	userRelService, serviceErr := service.NewUserRelService(database.Db)
	if serviceErr != nil {
		log.Error("NewUserRelService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userRelService", "Error happened while creating userRelService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[GetFollowersHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	followers, err := userRelService.GetFollowers(currentUser.UserID)
	if err != nil {
		log.Error("[GetFollowersHandle.userRelService.GetFollowers] %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/getFollowers", "Error happened while reading followers!"))
	}

	return c.JSON(followers)
}

// GetFollowingHandle handle get auth user following
func GetFollowingHandle(c *fiber.Ctx) error {

	// Create service
	userRelService, serviceErr := service.NewUserRelService(database.Db)
	if serviceErr != nil {
		log.Error("NewUserRelService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userRelService", "Error happened while creating userRelService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[GetFollowingHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	following, err := userRelService.GetFollowing(currentUser.UserID)
	if err != nil {
		log.Error("[GetFollowingHandle.userRelService.GetFollowing] %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/getFollowing", "Error happened while reading following!"))
	}

	return c.JSON(following)
}
