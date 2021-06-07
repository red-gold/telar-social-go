package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/votes/database"
	domain "github.com/red-gold/ts-serverless/micros/votes/dto"
	models "github.com/red-gold/ts-serverless/micros/votes/models"
	service "github.com/red-gold/ts-serverless/micros/votes/services"
)

// UpdateVoteHandle handle create a new vote
func UpdateVoteHandle(c *fiber.Ctx) error {

	// Create the model object
	var model models.VoteModel
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse VoteModel Error %s", err.Error())
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
		log.Error("[UpdateVoteHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	updatedVote := &domain.Vote{
		ObjectId:         model.ObjectId,
		PostId:           model.PostId,
		OwnerDisplayName: currentUser.DisplayName,
		OwnerAvatar:      currentUser.Avatar,
		CreatedDate:      model.CreatedDate,
	}

	if err := voteService.UpdateVoteById(updatedVote); err != nil {
		errorMessage := fmt.Sprintf("Update Vote Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updateVote", "Error happened while update Vote!"))
	}

	return c.SendStatus(http.StatusOK)
}
