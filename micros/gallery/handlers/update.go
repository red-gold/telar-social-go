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

// UpdateMediaHandle handle create a new media
func UpdateMediaHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.MediaModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse MediaModel Error %s", err.Error())
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
		log.Error("[UpdateMediaHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	updatedMedia := &domain.Media{
		ObjectId:       model.ObjectId,
		DeletedDate:    0,
		CreatedDate:    model.CreatedDate,
		Thumbnail:      model.Thumbnail,
		URL:            model.URL,
		FullPath:       model.FullPath,
		Caption:        model.Caption,
		FileName:       model.FileName,
		Directory:      model.Directory,
		OwnerUserId:    currentUser.UserID,
		LastUpdated:    utils.UTCNowUnix(),
		AlbumId:        model.AlbumId,
		Width:          model.Width,
		Height:         model.Height,
		Meta:           model.Meta,
		AccessUserList: model.AccessUserList,
		Permission:     model.Permission,
		Deleted:        model.Deleted,
	}

	if err := mediaService.UpdateMediaById(updatedMedia); err != nil {
		errorMessage := fmt.Sprintf("Update Media Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updateMedia", "Error happened while update media!"))
	}

	return c.SendStatus(http.StatusOK)

}
