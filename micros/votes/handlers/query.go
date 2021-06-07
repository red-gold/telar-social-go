package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/pkg/parser"
	utils "github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/votes/database"
	models "github.com/red-gold/ts-serverless/micros/votes/models"
	service "github.com/red-gold/ts-serverless/micros/votes/services"
)

type VoteQueryModel struct {
	Page   int64     `query:"page"`
	PostId uuid.UUID `query:"postId"`
}

// GetVotesByPostIdHandle handle query on vote
func GetVotesByPostIdHandle(c *fiber.Ctx) error {

	// Create service
	voteService, serviceErr := service.NewVoteService(database.Db)
	if serviceErr != nil {
		log.Error("NewVoteService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/voteService", "Error happened while creating voteService!"))
	}

	query := new(VoteQueryModel)

	if err := parser.QueryParser(c, query); err != nil {
		log.Error("[GetVotesByPostIdHandle] QueryParser %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryParser", "Error happened while parsing query!"))
	}

	if query.PostId == uuid.Nil {
		errorMessage := fmt.Sprintf("Post id can not be empty.")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("postIdRequired", errorMessage))
	}

	voteList, err := voteService.GetVoteByPostId(&query.PostId, "created_date", query.Page)
	if err != nil {
		log.Error("[GetVotesByPostIdHandle.voteService.GetVoteByPostId] %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/getVoteByPostId", "Error happened while query vote!"))
	}

	return c.JSON(voteList)
}

// GetVoteHandle handle get a vote
func GetVoteHandle(c *fiber.Ctx) error {

	// Create service
	voteService, serviceErr := service.NewVoteService(database.Db)
	if serviceErr != nil {
		log.Error("NewVoteService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/voteService", "Error happened while creating voteService!"))
	}
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

	foundVote, err := voteService.FindById(voteUUID)
	if err != nil {
		errorMessage := fmt.Sprintf("Find Vote %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/findVote", "Error happened while find Vote!"))
	}

	voteModel := models.VoteModel{
		ObjectId:         foundVote.ObjectId,
		OwnerUserId:      foundVote.OwnerUserId,
		PostId:           foundVote.PostId,
		OwnerDisplayName: foundVote.OwnerDisplayName,
		OwnerAvatar:      foundVote.OwnerAvatar,
		CreatedDate:      foundVote.CreatedDate,
		TypeId:           foundVote.TypeId,
	}

	return c.JSON(voteModel)
}
