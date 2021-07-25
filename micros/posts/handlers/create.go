package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/posts/database"
	domain "github.com/red-gold/ts-serverless/micros/posts/dto"
	models "github.com/red-gold/ts-serverless/micros/posts/models"
	service "github.com/red-gold/ts-serverless/micros/posts/services"
)

// CreatePostHandle handle create a new post
func CreatePostHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.CreatePostModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse CreatePostModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}
	var newAlbum *domain.PostAlbum = nil

	if len(model.Album.Photos) > 0 {
		newAlbum = &domain.PostAlbum{
			Count:   model.Album.Count,
			Cover:   model.Album.Cover,
			CoverId: model.Album.CoverId,
			Photos:  model.Album.Photos,
			Title:   model.Album.Title,
		}
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[CreatePostHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	newPost := &domain.Post{
		ObjectId:         model.ObjectId,
		PostTypeId:       model.PostTypeId,
		OwnerUserId:      currentUser.UserID,
		Score:            model.Score,
		Votes:            make(map[string]string),
		ViewCount:        model.ViewCount,
		Body:             model.Body,
		OwnerDisplayName: currentUser.DisplayName,
		OwnerAvatar:      currentUser.Avatar,
		URLKey:           generatPostURLKey(currentUser.SocialName, model.Body, model.ObjectId.String()),
		Tags:             model.Tags,
		CommentCounter:   model.CommentCounter,
		Image:            model.Image,
		ImageFullPath:    model.ImageFullPath,
		Video:            model.Video,
		Thumbnail:        model.Thumbnail,
		Album:            newAlbum,
		DisableComments:  model.DisableComments,
		DisableSharing:   model.DisableSharing,
		Deleted:          model.Deleted,
		DeletedDate:      model.DeletedDate,
		CreatedDate:      utils.UTCNowUnix(),
		LastUpdated:      model.LastUpdated,
		AccessUserList:   model.AccessUserList,
		Permission:       model.Permission,
		Version:          model.Version,
	}

	if err := postService.SavePost(newPost); err != nil {
		errorMessage := fmt.Sprintf("Save new post error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/savePost", "Error happened while save post!"))
	}

	return c.JSON(fiber.Map{
		"objectId": newPost.ObjectId.String(),
	})

}

// InitPostIndexHandle handle create a new post
func InitPostIndexHandle(c *fiber.Ctx) error {

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}

	postIndexMap := make(map[string]interface{})
	postIndexMap["body"] = "text"
	postIndexMap["objectId"] = 1
	if err := postService.CreatePostIndex(postIndexMap); err != nil {
		errorMessage := fmt.Sprintf("Create post index Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/createPostIndex", "Error happened while creating post index!"))
	}

	return c.SendStatus(http.StatusOK)

}
