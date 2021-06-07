package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/comments/database"
	domain "github.com/red-gold/ts-serverless/micros/comments/dto"
	models "github.com/red-gold/ts-serverless/micros/comments/models"
	service "github.com/red-gold/ts-serverless/micros/comments/services"
)

// UpdateCommentHandle handle create a new comment
func UpdateCommentHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.CommentModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse CommentModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	commentService, serviceErr := service.NewCommentService(database.Db)
	if serviceErr != nil {
		log.Error("NewCommentService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/commentService", "Error happened while creating commentService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[UpdateCommentHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	updatedComment := &domain.Comment{
		ObjectId:         model.ObjectId,
		OwnerUserId:      currentUser.UserID,
		PostId:           model.PostId,
		Score:            model.Score,
		Text:             model.Text,
		OwnerDisplayName: currentUser.DisplayName,
		OwnerAvatar:      currentUser.Avatar,
		Deleted:          model.Deleted,
		DeletedDate:      model.DeletedDate,
		CreatedDate:      model.CreatedDate,
		LastUpdated:      model.LastUpdated,
	}

	if err := commentService.UpdateCommentById(updatedComment); err != nil {
		errorMessage := fmt.Sprintf("Update Comment Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updateComment", "Error happened while update comment!"))
	}

	return c.SendStatus(http.StatusOK)

}

// UpdateCommentProfileHandle handle create a new post
func UpdateCommentProfileHandle(c *fiber.Ctx) error {

	// Create service
	postService, serviceErr := service.NewCommentService(database.Db)
	if serviceErr != nil {
		log.Error("NewCommentService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/commentService", "Error happened while creating commentService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[UpdateCommentProfileHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	postService.UpdateCommentProfile(currentUser.UserID, currentUser.DisplayName, currentUser.Avatar)

	return c.JSON(http.StatusOK)

}
