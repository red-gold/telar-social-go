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
	"github.com/red-gold/ts-serverless/micros/gallery/database"
	models "github.com/red-gold/ts-serverless/micros/gallery/models"
	service "github.com/red-gold/ts-serverless/micros/gallery/services"
)

type MediaQueryModel struct {
	Search string    `query:"search"`
	Page   int64     `query:"page"`
	Owner  uuid.UUID `query:"owner"`
	Type   int       `query:"type"`
}

type AlbumQueryModel struct {
	Page  int64     `query:"page"`
	Limit int64     `query:"limit"`
	Album uuid.UUID `query:"album"`
}

// QueryMediaHandle handle query on media
func QueryMediaHandle(c *fiber.Ctx) error {

	// Create service
	mediaService, serviceErr := service.NewMediaService(database.Db)
	if serviceErr != nil {
		log.Error("NewMediaService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/mediaService", "Error happened while creating mediaService!"))
	}

	query := new(MediaQueryModel)

	if err := parser.QueryParser(c, query); err != nil {
		log.Error("[QueryMediaHandle] QueryParser %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryParser", "Error happened while parsing query!"))
	}

	mediaList, err := mediaService.QueryMedia(query.Search, &query.Owner, &query.Type, "created_date", query.Page)
	if err != nil {
		log.Error("[QueryMediaHandle.mediaService.QueryMedia] %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/queryMedia", "Error happened while query media!"))

	}

	return c.JSON(mediaList)

}

// QueryAlbumHandle handle query on media
func QueryAlbumHandle(c *fiber.Ctx) error {

	// Create service
	mediaService, serviceErr := service.NewMediaService(database.Db)
	if serviceErr != nil {
		log.Error("NewMediaService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/mediaService", "Error happened while creating mediaService!"))
	}

	query := new(AlbumQueryModel)

	if err := parser.QueryParser(c, query); err != nil {
		log.Error("[QueryAlbumHandle] QueryParser %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryParser", "Error happened while parsing query!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[QueryAlbumHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	mediaList, err := mediaService.QueryAlbum(currentUser.UserID, &query.Album, query.Page, query.Limit, "created_date")
	if err != nil {
		log.Error("[QueryAlbumHandle.mediaService.QueryAlbum] %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/queryMedia", "Error happened while query media!"))
	}

	return c.JSON(mediaList)

}

// GetMediaHandle handle get a media
func GetMediaHandle(c *fiber.Ctx) error {

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

	foundMedia, err := mediaService.FindById(mediaUUID)
	if err != nil {
		log.Error("[GetMediaHandle.mediaService.FindById] %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/queryMedia", "Error happened while query media!"))
	}

	mediaModel := models.MediaModel{
		ObjectId:       foundMedia.ObjectId,
		DeletedDate:    foundMedia.DeletedDate,
		CreatedDate:    foundMedia.CreatedDate,
		Thumbnail:      foundMedia.Thumbnail,
		URL:            foundMedia.URL,
		FullPath:       foundMedia.FullPath,
		Caption:        foundMedia.Caption,
		FileName:       foundMedia.FileName,
		Directory:      foundMedia.Directory,
		OwnerUserId:    foundMedia.OwnerUserId,
		LastUpdated:    foundMedia.LastUpdated,
		AlbumId:        foundMedia.AlbumId,
		Width:          foundMedia.Width,
		Height:         foundMedia.Height,
		Meta:           foundMedia.Meta,
		AccessUserList: foundMedia.AccessUserList,
		Permission:     foundMedia.Permission,
		Deleted:        foundMedia.Deleted,
	}

	return c.JSON(mediaModel)

}

// GetMediaByDirectoryHandle handle get media list by directory
func GetMediaByDirectoryHandle(c *fiber.Ctx) error {

	// params from /medias/dir/:dir
	dirName := c.Params("dir")
	if dirName == "" {
		errorMessage := fmt.Sprintf("Directory name is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("dirNameRequired", errorMessage))
	}

	// Create service
	mediaService, serviceErr := service.NewMediaService(database.Db)
	if serviceErr != nil {
		log.Error("NewMediaService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/mediaService", "Error happened while creating mediaService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[GetMediaByDirectoryHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	foundMediaList, err := mediaService.FindByDirectory(currentUser.UserID, dirName, 0, 0)
	if err != nil {
		log.Error("[GetMediaByDirectoryHandle.mediaService.FindByDirectory] %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/queryMedia", "Error happened while query media!"))
	}

	return c.JSON(foundMediaList)

}
