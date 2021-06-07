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
	service "github.com/red-gold/ts-serverless/micros/posts/services"
)

// DeletePostHandle handle delete a post
func DeletePostHandle(c *fiber.Ctx) error {

	// params from /posts/:postId
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
		log.Error("[DeletePostHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := postService.DeletePostByOwner(currentUser.UserID, postUUID); err != nil {
		errorMessage := fmt.Sprintf("Delete Post Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/deletePost", "Error happened while deleting post!"))
	}

	return c.SendStatus(http.StatusOK)

}
