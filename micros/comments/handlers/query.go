package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/pkg/parser"
	utils "github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/comments/database"
	models "github.com/red-gold/ts-serverless/micros/comments/models"
	service "github.com/red-gold/ts-serverless/micros/comments/services"
)

type CommentQueryModel struct {
	Search string    `query:"search"`
	Page   int64     `query:"page"`
	Owner  uuid.UUID `query:"owner"`
	Type   int       `query:"type"`
}

type CommentQueryByPostIdModel struct {
	Page   int64     `query:"page"`
	PostId uuid.UUID `query:"postId"`
}

// QueryCommentHandle handle query on comment
func QueryCommentHandle(c *fiber.Ctx) error {

	// Create service
	commentService, serviceErr := service.NewCommentService(database.Db)
	if serviceErr != nil {
		log.Error("NewCommentService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/commentService", "Error happened while creating commentService!"))
	}

	query := new(CommentQueryModel)

	if err := parser.QueryParser(c, query); err != nil {
		log.Error("[QueryCommentHandle] QueryParser %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryParser", "Error happened while parsing query!"))
	}

	commentList, err := commentService.QueryComment(query.Search, &query.Owner, &query.Type, "created_date", query.Page)
	if err != nil {
		log.Error("[QueryCommentHandle.commentService.QueryComment] %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/queryComment", "Error happened while query comment!"))
	}

	return c.JSON(commentList)

}

// GetCommentsByPostIdHandle handle query on comment
func GetCommentsByPostIdHandle(c *fiber.Ctx) error {

	// Create service
	commentService, serviceErr := service.NewCommentService(database.Db)
	if serviceErr != nil {
		log.Error("NewCommentService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/commentService", "Error happened while creating commentService!"))
	}

	query := new(CommentQueryByPostIdModel)

	if err := parser.QueryParser(c, query); err != nil {
		log.Error("[GetCommentsByPostIdHandle] QueryParser %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryParser", "Error happened while parsing query!"))
	}

	if query.PostId == uuid.Nil {
		errorMessage := fmt.Sprintf("Post id can not be empty.")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdRequired", errorMessage))

	}

	commentList, err := commentService.GetCommentByPostId(&query.PostId, "created_date", query.Page)
	if err != nil {
		log.Error("[GetCommentsByPostIdHandle.commentService.GetCommentByPostId] %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/queryComment", "Error happened while query comment!"))
	}

	return c.JSON(commentList)

}

// GetCommentHandle handle get a comment
func GetCommentHandle(c *fiber.Ctx) error {

	// Create service
	commentService, serviceErr := service.NewCommentService(database.Db)
	if serviceErr != nil {
		log.Error("NewCommentService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/commentService", "Error happened while creating commentService!"))
	}
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

	foundComment, err := commentService.FindById(commentUUID)
	if err != nil {
		errorMessage := fmt.Sprintf("Find Comment %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/findComment", "Error happened while find comment!"))
	}

	// No comment found
	if foundComment == nil {
		return c.SendStatus(http.StatusOK)
	}
	commentModel := models.CommentModel{
		ObjectId:         foundComment.ObjectId,
		OwnerUserId:      foundComment.OwnerUserId,
		PostId:           foundComment.PostId,
		Score:            foundComment.Score,
		Text:             foundComment.Text,
		OwnerDisplayName: foundComment.OwnerDisplayName,
		OwnerAvatar:      foundComment.OwnerAvatar,
		Deleted:          foundComment.Deleted,
		DeletedDate:      foundComment.DeletedDate,
		CreatedDate:      foundComment.CreatedDate,
		LastUpdated:      foundComment.LastUpdated,
	}

	return c.JSON(commentModel)

}
