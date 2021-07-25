package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/posts/database"
	models "github.com/red-gold/ts-serverless/micros/posts/models"
	service "github.com/red-gold/ts-serverless/micros/posts/services"
)

// UpdatePostHandle handle create a new post
func UpdatePostHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.PostModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse PostModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}
	var updatedAlbum *models.PostAlbumModel = nil

	updatedAlbum = &models.PostAlbumModel{
		Count:   model.Album.Count,
		Cover:   model.Album.Cover,
		CoverId: model.Album.CoverId,
		Photos:  model.Album.Photos,
		Title:   model.Album.Title,
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[UpdatePostHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	updatedPost := &models.PostUpdateModel{
		ObjectId:         model.ObjectId,
		PostTypeId:       model.PostTypeId,
		OwnerUserId:      currentUser.UserID,
		Score:            model.Score,
		Votes:            model.Votes,
		ViewCount:        model.ViewCount,
		Body:             model.Body,
		OwnerDisplayName: currentUser.DisplayName,
		OwnerAvatar:      currentUser.Avatar,
		Tags:             model.Tags,
		CommentCounter:   model.CommentCounter,
		Image:            model.Image,
		ImageFullPath:    model.ImageFullPath,
		Video:            model.Video,
		Thumbnail:        model.Thumbnail,
		Album:            updatedAlbum,
		DisableComments:  model.DisableComments,
		DisableSharing:   model.DisableSharing,
		LastUpdated:      utils.UTCNowUnix(),
		AccessUserList:   model.AccessUserList,
		Permission:       model.Permission,
		Version:          model.Version,
	}

	if err := postService.UpdatePostById(updatedPost); err != nil {
		errorMessage := fmt.Sprintf("Update Post Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))
	}

	return c.SendStatus(http.StatusOK)

}

// IncrementScoreHandle handle create a new post
func IncrementScoreHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.ScoreModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse ScoreModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	if model.PostId == uuid.Nil {
		errorMessage := fmt.Sprintf("Post Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdRequired", errorMessage))
	}
	if model.Count == 0 {
		errorMessage := fmt.Sprintf("Count can not be zero!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("countIsZero", errorMessage))
	}

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[IncrementScoreHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if model.Count > 0 {
		err := postService.IncrementScoreCount(model.PostId, currentUser.UserID, currentUser.Avatar)
		if err != nil {
			errorMessage := fmt.Sprintf("[IncrementScoreCount] Update Post Error %s", err.Error())
			log.Error(errorMessage)
			return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))

		}
	} else if model.Count < 0 {
		err := postService.DecrementScoreCount(model.PostId, currentUser.UserID)
		if err != nil {
			errorMessage := fmt.Sprintf("[DecrementScoreCount] Update Post Error %s", err.Error())
			log.Error(errorMessage)
			return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))

		}
	}

	return c.SendStatus(http.StatusOK)

}

// IncrementCommentHandle handle create a new post
func IncrementCommentHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.CommentCountModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse CommentCountModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	if model.PostId == uuid.Nil {
		errorMessage := fmt.Sprintf("Post Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdRequired", errorMessage))
	}
	if model.Count == 0 {
		errorMessage := fmt.Sprintf("Count can not be zero!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("countIsZero", errorMessage))
	}

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}

	if model.Count > 0 {
		err := postService.IncrementCommentCount(model.PostId)
		if err != nil {
			errorMessage := fmt.Sprintf("[IncrementCommentCount] Update Post Error %s", err.Error())
			log.Error(errorMessage)
			return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))

		}
	} else if model.Count < 0 {
		err := postService.DecerementCommentCount(model.PostId)
		if err != nil {
			errorMessage := fmt.Sprintf("[DecerementCommentCount] Update Post Error %s", err.Error())
			log.Error(errorMessage)
			return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))

		}
	}
	return c.SendStatus(http.StatusOK)

}

// DisableCommentHandle disble post's commnet
func DisableCommentHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.DisableCommentModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse DisableCommentModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	if model.PostId == uuid.Nil {
		errorMessage := fmt.Sprintf("Post Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdRequired", errorMessage))
	}

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[DisableCommentHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	err := postService.DisableCommnet(currentUser.UserID, model.PostId, model.Status)
	if err != nil {
		errorMessage := fmt.Sprintf("[DisableCommnet] Update Post Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))

	}
	return c.SendStatus(http.StatusOK)

}

// DisableSharingHandle disble post's sharing
func DisableSharingHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.DisableSharingModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse DisableSharingModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	if model.PostId == uuid.Nil {
		errorMessage := fmt.Sprintf("Post Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdRequired", errorMessage))
	}

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[DisableSharingHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	err := postService.DisableSharing(currentUser.UserID, model.PostId, model.Status)
	if err != nil {
		errorMessage := fmt.Sprintf("[DisableSharing] Update Post Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))

	}
	return c.SendStatus(http.StatusOK)

}

// Deprecated: UpdatePostProfileHandle
func UpdatePostProfileHandle(c *fiber.Ctx) error {

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[UpdatePostProfileHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	err := postService.UpdatePostProfile(currentUser.UserID, currentUser.DisplayName, currentUser.Avatar)
	if err != nil {
		errorMessage := fmt.Sprintf("[UpdatePostProfile] Update Post Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))

	}

	return c.SendStatus(http.StatusOK)

}
