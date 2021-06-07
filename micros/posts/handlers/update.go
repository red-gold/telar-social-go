package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/constants"
	"github.com/red-gold/ts-serverless/micros/posts/database"
	domain "github.com/red-gold/ts-serverless/micros/posts/dto"
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
	var updatedAlbum *domain.PostAlbum = nil
	if model.PostTypeId == constants.PostConstAlbum.Parse() || model.PostTypeId == constants.PostConstPhotoGallery.Parse() {
		updatedAlbum = &domain.PostAlbum{
			Count:   model.Album.Count,
			Cover:   model.Album.Cover,
			CoverId: model.Album.CoverId,
			Photos:  model.Album.Photos,
			Title:   model.Album.Title,
		}
	}
	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[UpdatePostHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	updatedPost := &domain.Post{
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
		Deleted:          model.Deleted,
		DeletedDate:      model.DeletedDate,
		CreatedDate:      model.CreatedDate,
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

	// params from /posts/score/+1/:postId
	postId := c.Params("postId")
	if postId == "" {
		errorMessage := fmt.Sprintf("Post Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdRequired", errorMessage))
	}

	postUUID, uuidErr := uuid.FromString(postId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdIsNotValid", "Post id is not valid!"))
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

	err := postService.IncrementScoreCount(postUUID, currentUser.UserID)
	if err != nil {
		errorMessage := fmt.Sprintf("[IncrementScoreCount] Update Post Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))

	}

	return c.SendStatus(http.StatusOK)

}

// DecrementScoreHandle handle create a new post
func DecrementScoreHandle(c *fiber.Ctx) error {

	// params from /posts/score/-1/:postId
	postId := c.Params("postId")
	if postId == "" {
		errorMessage := fmt.Sprintf("Post Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdRequired", errorMessage))
	}

	postUUID, uuidErr := uuid.FromString(postId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdIsNotValid", "Post id is not valid!"))
	}

	fmt.Printf("\n Post UUID: %s", postUUID)
	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[DecrementScoreHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	err := postService.DecrementScoreCount(postUUID, currentUser.UserID)
	if err != nil {
		errorMessage := fmt.Sprintf("[DecrementScoreCount] Update Post Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))

	}
	return c.SendStatus(http.StatusOK)

}

// IncrementCommentHandle handle create a new post
func IncrementCommentHandle(c *fiber.Ctx) error {

	// params from /comment/+1/:postId
	postId := c.Params("postId")
	if postId == "" {
		errorMessage := fmt.Sprintf("Post Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdRequired", errorMessage))
	}

	postUUID, uuidErr := uuid.FromString(postId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdIsNotValid", "Post id is not valid!"))
	}

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}

	err := postService.IncrementCommentCount(postUUID)
	if err != nil {
		errorMessage := fmt.Sprintf("[IncrementCommentCount] Update Post Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))

	}

	return c.SendStatus(http.StatusOK)

}

// DecrementCommentHandle handle create a new post
func DecrementCommentHandle(c *fiber.Ctx) error {

	// params from /comment/-1/:postId
	postId := c.Params("postId")
	if postId == "" {
		errorMessage := fmt.Sprintf("Post Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdRequired", errorMessage))
	}

	postUUID, uuidErr := uuid.FromString(postId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdIsNotValid", "Post id is not valid!"))
	}

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}

	err := postService.DecerementCommentCount(postUUID)
	if err != nil {
		errorMessage := fmt.Sprintf("[DecerementCommentCount] Update Post Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))

	}

	return c.SendStatus(http.StatusOK)

}

// Deprecated: UpdatePostProfileHandle handle create a new post
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
