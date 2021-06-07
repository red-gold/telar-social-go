package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/comments/database"
	service "github.com/red-gold/ts-serverless/micros/comments/services"
)

// DeleteCommentHandle handle delete a Comment
func DeleteCommentHandle(c *fiber.Ctx) error {

	// params from /comments/id/:commentId/post/:postId
	commentId := c.Params("commentId")
	if commentId == "" {
		errorMessage := fmt.Sprintf("Comment Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("commentIdRequired", errorMessage))
	}

	commentUUID, uuidErr := uuid.FromString(commentId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("commentIdIsNotValid", "Comment id is not valid!"))
	}

	postId := c.Params("postId")
	if postId == "" {
		errorMessage := fmt.Sprintf("Post Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdRequired", errorMessage))
	}

	// Create service
	commentService, serviceErr := service.NewCommentService(database.Db)
	if serviceErr != nil {
		log.Error("NewCommentService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/commentService", "Error happened while creating commentService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[DeleteCommentHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := commentService.DeleteCommentByOwner(currentUser.UserID, commentUUID); err != nil {
		errorMessage := fmt.Sprintf("Delete Comment Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/deleteComment", "Error happened while delete comment!"))
	}

	// Create user headers for http request
	userHeaders := make(map[string][]string)
	userHeaders["uid"] = []string{currentUser.UserID.String()}
	userHeaders["email"] = []string{currentUser.Username}
	userHeaders["avatar"] = []string{currentUser.Avatar}
	userHeaders["displayName"] = []string{currentUser.DisplayName}
	userHeaders["role"] = []string{currentUser.SystemRole}

	postIndexURL := fmt.Sprintf("/posts/comment/-1/%s", postId)
	_, commentDecreaseRes := functionCall(http.MethodPut, []byte(""), postIndexURL, userHeaders)

	if commentDecreaseRes != nil {
		log.Error("Cannot save vote on post! error: %s", commentDecreaseRes.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/decreasePostCommentCount", "Error happened while decrease post comment count!"))
	}

	return c.SendStatus(http.StatusOK)

}

// DeleteCommentByPostIdHandle handle delete a Comment but postId
func DeleteCommentByPostIdHandle(c *fiber.Ctx) error {

	// params from /Comments/post/:postId
	postId := c.Params("postId")
	if postId == "" {
		errorMessage := fmt.Sprintf("Post Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdRequired", errorMessage))
	}

	PostUUID, uuidErr := uuid.FromString(postId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdIsNotValid", "Post id is not valid!"))
	}

	// Create service
	commentService, serviceErr := service.NewCommentService(database.Db)
	if serviceErr != nil {
		log.Error("NewCommentService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/commentService", "Error happened while creating commentService!"))

	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[DeleteCommentByPostIdHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := commentService.DeleteCommentsByPostId(currentUser.UserID, PostUUID); err != nil {
		errorMessage := fmt.Sprintf("Delete Comment Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/deleteComment", "Error happened while delete comment!"))
	}

	return c.SendStatus(http.StatusOK)

}
