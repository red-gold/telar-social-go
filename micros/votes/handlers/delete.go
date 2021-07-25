package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/votes/database"
	service "github.com/red-gold/ts-serverless/micros/votes/services"
)

// DeleteVoteHandle handle delete a Vote
func DeleteVoteHandle(c *fiber.Ctx) error {

	// params from /votes/id/:voteId
	voteId := c.Params("voteId")
	if voteId == "" {
		errorMessage := fmt.Sprintf("Vote Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("voteIdRequired", errorMessage))
	}

	voteUUID, uuidErr := uuid.FromString(voteId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("voteIdIsNotValid", "Vote id is not valid!"))

	}
	fmt.Printf("\n Vote UUID: %s", voteUUID)
	// Create service
	voteService, serviceErr := service.NewVoteService(database.Db)
	if serviceErr != nil {
		log.Error("NewVoteService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/voteService", "Error happened while creating voteService!"))

	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[DeleteVoteHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := voteService.DeleteVoteByOwner(currentUser.UserID, voteUUID); err != nil {
		errorMessage := fmt.Sprintf("Delete Vote Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/deleteVote", "Error happened while delete Vote!"))

	}

	return c.SendStatus(http.StatusOK)
}

// DeleteVoteByPostIdHandle handle delete a Vote but postId
func DeleteVoteByPostIdHandle(c *fiber.Ctx) error {

	// params from /Votes/post/:postId
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
	voteService, serviceErr := service.NewVoteService(database.Db)
	if serviceErr != nil {
		log.Error("NewVoteService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/voteService", "Error happened while creating voteService!"))

	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[DeleteVoteByPostIdHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := voteService.DeleteVotesByPostId(currentUser.UserID, PostUUID); err != nil {
		errorMessage := fmt.Sprintf("Delete Vote Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/deleteVote", "Error happened while delete Vote!"))
	}

	// Create user headers for http request
	userHeaders := make(map[string][]string)
	userHeaders["uid"] = []string{currentUser.UserID.String()}
	userHeaders["email"] = []string{currentUser.Username}
	userHeaders["avatar"] = []string{currentUser.Avatar}
	userHeaders["displayName"] = []string{currentUser.DisplayName}
	userHeaders["role"] = []string{currentUser.SystemRole}

	fullURL := "/posts/score"
	payload, err := json.Marshal(fiber.Map{
		"postId": postId,
		"count":  -1,
	})
	if err != nil {
		messageError := fmt.Sprintf("Can not parse score payload: %s", err.Error())
		log.Error(messageError)
	}

	_, functionErr := functionCall(http.MethodPut, payload, fullURL, userHeaders)
	if functionErr != nil {
		log.Error("[DeleteVoteByPostIdHandle.functionCall] %s - %s", fullURL, functionErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postDecreaseScore", "Error happened while delete Vote!"))
	}

	return c.SendStatus(http.StatusOK)
}
