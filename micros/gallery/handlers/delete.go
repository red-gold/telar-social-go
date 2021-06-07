package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/gallery/database"
	service "github.com/red-gold/ts-serverless/micros/gallery/services"
)

// DeleteMediaHandle handle delete a media
func DeleteMediaHandle(c *fiber.Ctx) error {

	// params from /medias/id/:mediaId
	mediaId := c.Params("mediaId")
	if mediaId == "" {
		errorMessage := fmt.Sprintf("Media Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("mediaIdRequired", errorMessage))
	}

	mediaUUID, uuidErr := uuid.FromString(mediaId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("mediaIdIsNotValid", "Media id is not valid!"))
	}

	// Create service
	mediaService, serviceErr := service.NewMediaService(database.Db)
	if serviceErr != nil {
		log.Error("NewMediaService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/mediaService", "Error happened while creating mediaService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[DeleteMediaHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := mediaService.DeleteMediaByOwner(currentUser.UserID, mediaUUID); err != nil {
		errorMessage := fmt.Sprintf("Delete Media Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/deleteMedia", "Error happened while delete media!"))
	}

	return c.SendStatus(http.StatusOK)

}

// DeleteDirectoryHandle handle delete a media
func DeleteDirectoryHandle(c *fiber.Ctx) error {

	// params from /medias/dir/:dir
	dirName := c.Params("dir")
	if dirName == "" {
		errorMessage := fmt.Sprintf("Directory name is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("directoryNameIsRequired", errorMessage))
	}

	// Create service
	mediaService, serviceErr := service.NewMediaService(database.Db)
	if serviceErr != nil {
		log.Error("NewMediaService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/mediaService", "Error happened while creating mediaService!"))

	}
	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[DeleteDirectoryHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := mediaService.DeleteMediaByDirectory(currentUser.UserID, dirName); err != nil {
		errorMessage := fmt.Sprintf("Delete Media Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/deleteMedia", "Error happened while delete media!"))
	}

	return c.SendStatus(http.StatusOK)

}
