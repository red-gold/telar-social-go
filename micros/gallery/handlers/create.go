package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/gallery/database"
	domain "github.com/red-gold/ts-serverless/micros/gallery/dto"
	models "github.com/red-gold/ts-serverless/micros/gallery/models"
	service "github.com/red-gold/ts-serverless/micros/gallery/services"
)

// CreateMediaHandle handle create a new media
func CreateMediaHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.CreateMediaModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse CreateMediaModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	mediaService, serviceErr := service.NewMediaService(database.Db)
	if serviceErr != nil {
		log.Error("NewMediaService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/mediaService", "Error happened while creating mediaService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[CreateMediaHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	newMedia := &domain.Media{
		ObjectId:       model.ObjectId,
		DeletedDate:    0,
		CreatedDate:    utils.UTCNowUnix(),
		Thumbnail:      model.Thumbnail,
		URL:            model.URL,
		FullPath:       model.FullPath,
		Caption:        model.Caption,
		FileName:       model.FileName,
		Directory:      model.Directory,
		OwnerUserId:    currentUser.UserID,
		LastUpdated:    0,
		AlbumId:        model.AlbumId,
		Width:          model.Width,
		Height:         model.Height,
		Meta:           model.Meta,
		AccessUserList: model.AccessUserList,
		Permission:     model.Permission,
		Deleted:        false,
	}

	if err := mediaService.SaveMedia(newMedia); err != nil {
		errorMessage := fmt.Sprintf("Save Media Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveMedia", "Error happened while saving media!"))
	}

	return c.JSON(fiber.Map{
		"objectId": newMedia.ObjectId.String(),
	})

}

// CreateMediaListHandle handle create a new media
func CreateMediaListHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.CreateMediaListModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse CreateMediaListModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	mediaService, serviceErr := service.NewMediaService(database.Db)
	if serviceErr != nil {
		log.Error("NewMediaService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/mediaService", "Error happened while creating mediaService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[CreateMediaListHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	var mediaList []domain.Media
	for _, media := range model.List {

		newMedia := domain.Media{
			ObjectId:       media.ObjectId,
			DeletedDate:    0,
			CreatedDate:    utils.UTCNowUnix(),
			Thumbnail:      media.Thumbnail,
			URL:            media.URL,
			FullPath:       media.FullPath,
			Caption:        media.Caption,
			FileName:       media.FileName,
			Directory:      media.Directory,
			OwnerUserId:    currentUser.UserID,
			LastUpdated:    0,
			AlbumId:        media.AlbumId,
			Width:          media.Width,
			Height:         media.Height,
			Meta:           media.Meta,
			AccessUserList: media.AccessUserList,
			Permission:     media.Permission,
			Deleted:        false,
		}
		mediaList = append(mediaList, newMedia)
	}

	if err := mediaService.SaveManyMedia(mediaList); err != nil {
		errorMessage := fmt.Sprintf("Save Media Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveMedia", "Error happened while saving media!"))
	}

	return c.SendStatus(http.StatusOK)

}
