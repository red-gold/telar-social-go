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
	"github.com/red-gold/ts-serverless/micros/comments/database"
	domain "github.com/red-gold/ts-serverless/micros/comments/dto"
	models "github.com/red-gold/ts-serverless/micros/comments/models"
	service "github.com/red-gold/ts-serverless/micros/comments/services"
)

type PostModelNotification struct {
	ObjectId         uuid.UUID `json:"objectId"`
	OwnerUserId      uuid.UUID `json:"ownerUserId"`
	OwnerDisplayName string    `json:"ownerDisplayName"`
	OwnerAvatar      string    `json:"ownerAvatar"`
}

// CreateCommentHandle handle create a new comment
func CreateCommentHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.CreateCommentModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse CreateCommentModel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	if model.Text == "" {
		errorMessage := fmt.Sprintf("Comment text is required")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("textIsRequired", errorMessage))
	}

	if model.PostId == uuid.Nil {
		errorMessage := fmt.Sprintf("Comment postId is required")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdIsRequired", errorMessage))
	}

	// Create service
	commentService, serviceErr := service.NewCommentService(database.Db)
	if serviceErr != nil {
		log.Error("NewCommentService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/commentService", "Error happened while creating commentService!"))
	}

	currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	if !ok {
		log.Error("[CreateCommentHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	newComment := &domain.Comment{
		OwnerUserId:      currentUser.UserID,
		PostId:           model.PostId,
		Score:            0,
		Text:             model.Text,
		OwnerDisplayName: currentUser.DisplayName,
		OwnerAvatar:      currentUser.Avatar,
		Deleted:          false,
		DeletedDate:      0,
		CreatedDate:      utils.UTCNowUnix(),
		LastUpdated:      0,
	}

	if err := commentService.SaveComment(newComment); err != nil {
		errorMessage := fmt.Sprintf("Save Comment Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveComment", "Error happened while saving comment!"))
	}

	// Create user headers for http request
	userHeaders := make(map[string][]string)
	userHeaders["uid"] = []string{currentUser.UserID.String()}
	userHeaders["email"] = []string{currentUser.Username}
	userHeaders["avatar"] = []string{currentUser.Avatar}
	userHeaders["displayName"] = []string{currentUser.DisplayName}
	userHeaders["role"] = []string{currentUser.SystemRole}

	// Create request to increase comment counter on post
	go func() {

		postCommentURL := fmt.Sprintf("/posts/comment/+1/%s", model.PostId)
		_, postCommentErr := functionCall(http.MethodPut, []byte(""), postCommentURL, userHeaders)

		if postCommentErr != nil {
			messageError := fmt.Sprintf("Cannot save comment count on post! error: %s", postCommentErr.Error())
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
			// Should not send notification if the owner of the comment is same as owner of post
			return
		}
		URL := fmt.Sprintf("/%s/posts/%s", currentUser.UserID, model.PostId)
		notificationModel := &notificationsModels.CreateNotificationModel{
			OwnerUserId:          currentUser.UserID,
			OwnerDisplayName:     currentUser.DisplayName,
			OwnerAvatar:          currentUser.Avatar,
			Description:          fmt.Sprintf("%s commented on your post.", currentUser.DisplayName),
			URL:                  URL,
			NotifyRecieverUserId: post.OwnerUserId,
			TargetId:             model.PostId,
			IsSeen:               false,
			Type:                 "comment",
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
		"objectId": newComment.ObjectId.String(),
	})

}
