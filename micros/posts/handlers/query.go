package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/pkg/parser"
	utils "github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/micros/posts/database"
	models "github.com/red-gold/ts-serverless/micros/posts/models"
	service "github.com/red-gold/ts-serverless/micros/posts/services"
)

type PostQueryModel struct {
	Search string      `query:"search"`
	Page   int64       `query:"page"`
	Owner  []uuid.UUID `query:"owner"`
	Type   int         `query:"type"`
}

// QueryPostHandle handle query on post
func QueryPostHandle(c *fiber.Ctx) error {

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}

	query := new(PostQueryModel)

	if err := parser.QueryParser(c, query); err != nil {
		log.Error("[QueryPostHandle] QueryParser %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryParser", "Error happened while parsing query!"))
	}

	postList, err := postService.QueryPostIncludeUser(query.Search, query.Owner, query.Type, "created_date", query.Page)
	if err != nil {
		log.Error("[QueryPostHandle.postService.QueryPostIncludeUser] %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/queryPost", "Error happened while query post!"))
	}

	return c.JSON(postList)

}

// GetPostHandle handle get a post
func GetPostHandle(c *fiber.Ctx) error {

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}
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

	foundPost, err := postService.FindById(postUUID)
	if err != nil {
		log.Error("[GetPostHandle.postService.FindById] %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/queryPost", "Error happened while query post!"))
	}

	postModel := models.PostModel{
		ObjectId:         foundPost.ObjectId,
		PostTypeId:       foundPost.PostTypeId,
		OwnerUserId:      foundPost.OwnerUserId,
		Score:            foundPost.Score,
		Votes:            foundPost.Votes,
		ViewCount:        foundPost.ViewCount,
		Body:             foundPost.Body,
		OwnerDisplayName: foundPost.OwnerDisplayName,
		OwnerAvatar:      foundPost.OwnerAvatar,
		URLKey:           foundPost.URLKey,
		Album:            models.PostAlbumModel{Photos: []string{}},
		Tags:             foundPost.Tags,
		CommentCounter:   foundPost.CommentCounter,
		Image:            foundPost.Image,
		ImageFullPath:    foundPost.ImageFullPath,
		Video:            foundPost.Video,
		Thumbnail:        foundPost.Thumbnail,
		DisableComments:  foundPost.DisableComments,
		DisableSharing:   foundPost.DisableSharing,
		Deleted:          foundPost.Deleted,
		DeletedDate:      foundPost.DeletedDate,
		CreatedDate:      foundPost.CreatedDate,
		LastUpdated:      foundPost.LastUpdated,
		AccessUserList:   foundPost.AccessUserList,
		Permission:       foundPost.Permission,
		Version:          foundPost.Version,
	}

	if foundPost.Album != nil && len(foundPost.Album.Photos) > 0 {
		postModel.Album = models.PostAlbumModel{
			Count:   foundPost.Album.Count,
			Cover:   foundPost.Album.Cover,
			CoverId: foundPost.Album.CoverId,
			Photos:  foundPost.Album.Photos,
			Title:   foundPost.Album.Title,
		}
	}

	return c.JSON(postModel)

}

// GetPostByURLKeyHandle handle get a post
func GetPostByURLKeyHandle(c *fiber.Ctx) error {

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}
	urlKey := c.Params("urlkey")
	if urlKey == "" {
		errorMessage := fmt.Sprintf("URL key is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("urlKeyRequired", errorMessage))
	}

	foundPost, err := postService.FindByURLKey(urlKey)
	if err != nil {
		log.Error("[GetPostHandle.postService.FindByURLKey] %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/queryPost", "Error happened while query post!"))
	}
	postModel := models.PostModel{
		ObjectId:         foundPost.ObjectId,
		PostTypeId:       foundPost.PostTypeId,
		OwnerUserId:      foundPost.OwnerUserId,
		Score:            foundPost.Score,
		Votes:            foundPost.Votes,
		ViewCount:        foundPost.ViewCount,
		Body:             foundPost.Body,
		OwnerDisplayName: foundPost.OwnerDisplayName,
		OwnerAvatar:      foundPost.OwnerAvatar,
		URLKey:           foundPost.URLKey,
		Album:            models.PostAlbumModel{Photos: []string{}},
		Tags:             foundPost.Tags,
		CommentCounter:   foundPost.CommentCounter,
		Image:            foundPost.Image,
		ImageFullPath:    foundPost.ImageFullPath,
		Video:            foundPost.Video,
		Thumbnail:        foundPost.Thumbnail,
		DisableComments:  foundPost.DisableComments,
		DisableSharing:   foundPost.DisableSharing,
		Deleted:          foundPost.Deleted,
		DeletedDate:      foundPost.DeletedDate,
		CreatedDate:      foundPost.CreatedDate,
		LastUpdated:      foundPost.LastUpdated,
		AccessUserList:   foundPost.AccessUserList,
		Permission:       foundPost.Permission,
		Version:          foundPost.Version,
	}
	log.Info("postModel %v", postModel)

	if foundPost.Album != nil && len(foundPost.Album.Photos) > 0 {
		postModel.Album = models.PostAlbumModel{
			Count:   foundPost.Album.Count,
			Cover:   foundPost.Album.Cover,
			CoverId: foundPost.Album.CoverId,
			Photos:  foundPost.Album.Photos,
			Title:   foundPost.Album.Title,
		}
	}

	return c.JSON(postModel)

}

// GeneratePostURLKeyHandle handle get post URL key
func GeneratePostURLKeyHandle(c *fiber.Ctx) error {

	// Create service
	postService, serviceErr := service.NewPostService(database.Db)
	if serviceErr != nil {
		log.Error("NewPostService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/postService", "Error happened while creating postService!"))
	}
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

	foundPost, err := postService.FindById(postUUID)
	if err != nil {
		log.Error("[GetPostHandle.postService.FindById] %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/queryPost", "Error happened while query post!"))
	}

	if foundPost.URLKey == "" {
		postOwnerProfile, err := getUserProfileByID(foundPost.OwnerUserId)

		if err != nil {
			log.Error("[GetPostHandle.getUserProfileByID] %s ", err.Error())
			return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/queryOwnerProfilePost", "Error happened while query owner profile post!"))
		}

		urlKey := generatPostURLKey(postOwnerProfile.SocialName, foundPost.Body, foundPost.ObjectId.String())
		err = postService.UpdatePostURLKey(foundPost.ObjectId, urlKey)
		if err != nil {
			log.Error("[GetPostHandle.postService.UpdatePostURLKey] %s ", err.Error())
			return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updatePost", "Error happened while updating post!"))
		}
		return c.JSON(fiber.Map{
			"urlKey": urlKey,
		})
	}

	return c.JSON(fiber.Map{
		"urlKey": foundPost.URLKey,
	})

}
