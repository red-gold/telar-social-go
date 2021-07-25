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
	domain "github.com/red-gold/ts-serverless/micros/votes/dto"
	models "github.com/red-gold/ts-serverless/micros/votes/models"
	service "github.com/red-gold/ts-serverless/micros/votes/services"
)

type PostModelNotification struct {
	ObjectId         uuid.UUID `json:"objectId"`
	OwnerUserId      uuid.UUID `json:"ownerUserId"`
	OwnerDisplayName string    `json:"ownerDisplayName"`
	OwnerAvatar      string    `json:"ownerAvatar"`
	URLKey           string    `json:"urlKey"`
}

// CreateVoteHandle handle create a new vote
func CreateVoteHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.CreateVoteModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse CreateVoteModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	voteService, serviceErr := service.NewVoteService(database.Db)
	if serviceErr != nil {
		log.Error("NewVoteService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/voteService", "Error happened while creating voteService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[CreateVoteHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	newVote := &domain.Vote{
		OwnerUserId:      currentUser.UserID,
		PostId:           model.PostId,
		OwnerDisplayName: currentUser.DisplayName,
		OwnerAvatar:      currentUser.Avatar,
		TypeId:           model.TypeId,
	}
	userInfoReq := getUserInfoReq(c)

	saveVoteChannel := voteService.SaveVote(newVote)
	readPostChannel := readPostAsync(model.PostId, userInfoReq)

	saveVoteResult, postResult := <-saveVoteChannel, <-readPostChannel
	if saveVoteResult.Error != nil {
		errorMessage := fmt.Sprintf("Save Vote Error %s", saveVoteResult.Error.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveVote", "Error happened while saving Vote!"))
	}

	if postResult.Error != nil {
		messageError := fmt.Sprintf("Cannot get the post! error: %s", postResult.Error.Error())
		fmt.Println(messageError)
	}

	// Create user headers for http request
	userHeaders := make(map[string][]string)
	userHeaders["uid"] = []string{currentUser.UserID.String()}
	userHeaders["email"] = []string{currentUser.Username}
	userHeaders["avatar"] = []string{currentUser.Avatar}
	userHeaders["displayName"] = []string{currentUser.DisplayName}
	userHeaders["role"] = []string{currentUser.SystemRole}

	go func() {
		fullURL := "/posts/score"
		payload, err := json.Marshal(fiber.Map{
			"postId": model.PostId,
			"count":  1,
		})
		if err != nil {
			messageError := fmt.Sprintf("Can not parse score payload: %s", err.Error())
			log.Error(messageError)
		}

		_, functionErr := functionCall(http.MethodPut, payload, fullURL, userHeaders)

		if functionErr != nil {
			messageError := fmt.Sprintf("Cannot save vote on post! error: %s", functionErr.Error())
			log.Error(messageError)
		}

	}()

	// Create notification request
	go func(currentUser types.UserContext) {

		var post PostModelNotification
		marshalErr := json.Unmarshal(postResult.Result, &post)
		if marshalErr != nil {
			messageError := fmt.Sprintf("Cannot unmarshal the post! error: %s", marshalErr.Error())
			fmt.Println(messageError)
		}

		if post.OwnerUserId == currentUser.UserID {
			// Should not send notification if the owner of the vote is same as owner of post
			return
		}

		URL := fmt.Sprintf("/posts/%s", post.URLKey)
		notificationModel := &models.NotificationModel{
			OwnerUserId:          currentUser.UserID,
			OwnerDisplayName:     currentUser.DisplayName,
			OwnerAvatar:          currentUser.Avatar,
			Title:                currentUser.DisplayName,
			Description:          fmt.Sprintf("%s like your post.", currentUser.DisplayName),
			URL:                  URL,
			NotifyRecieverUserId: post.OwnerUserId,
			TargetId:             model.PostId,
			IsSeen:               false,
			Type:                 "like",
		}
		notificationBytes, marshalErr := json.Marshal(notificationModel)
		if marshalErr != nil {
			fmt.Printf("Cannot marshal notification! error: %s", marshalErr.Error())

		}

		notificationURL := "/notifications"
		_, notificationIndexErr := functionCall(http.MethodPost, notificationBytes, notificationURL, userHeaders)
		if notificationIndexErr != nil {
			fmt.Printf("\nCannot save notification on follow user! error: %s\n", notificationIndexErr.Error())
		}

	}(currentUser)

	return c.JSON(fiber.Map{
		"objectId": newVote.ObjectId.String(),
	})
}
