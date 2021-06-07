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
	notificationsModels "github.com/red-gold/telar-web/micros/notifications/models"
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

	if err := voteService.SaveVote(newVote); err != nil {
		errorMessage := fmt.Sprintf("Save Vote Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveVote", "Error happened while saving Vote!"))
	}

	// Create user headers for http request
	userHeaders := make(map[string][]string)
	userHeaders["uid"] = []string{currentUser.UserID.String()}
	userHeaders["email"] = []string{currentUser.Username}
	userHeaders["avatar"] = []string{currentUser.Avatar}
	userHeaders["displayName"] = []string{currentUser.DisplayName}
	userHeaders["role"] = []string{currentUser.SystemRole}

	go func() {
		postURL := fmt.Sprintf("/posts/score/+1/%s", model.PostId)
		_, postErr := functionCall(http.MethodPut, []byte(""), postURL, userHeaders)

		if postErr != nil {
			messageError := fmt.Sprintf("Cannot save vote on post! error: %s", postErr.Error())
			fmt.Println(messageError)
		}

	}()

	// Create notification request
	go func(currentUser types.UserContext) {
		postURL := fmt.Sprintf("/posts/%s", model.PostId)
		postBody, postErr := functionCall(http.MethodGet, []byte(""), postURL, userHeaders)

		if postErr != nil {
			messageError := fmt.Sprintf("Cannot get the post! error: %s", postErr.Error())
			fmt.Println(messageError)
		}

		var post PostModelNotification
		marshalErr := json.Unmarshal(postBody, &post)
		if marshalErr != nil {
			messageError := fmt.Sprintf("Cannot unmarshal the post! error: %s", marshalErr.Error())
			fmt.Println(messageError)
		}

		if post.OwnerUserId == currentUser.UserID {
			// Should not send notification if the owner of the vote is same as owner of post
			return
		}

		URL := fmt.Sprintf("/%s/posts/%s", currentUser.UserID, model.PostId)
		notificationModel := &notificationsModels.CreateNotificationModel{
			OwnerUserId:          currentUser.UserID,
			OwnerDisplayName:     currentUser.DisplayName,
			OwnerAvatar:          currentUser.Avatar,
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
