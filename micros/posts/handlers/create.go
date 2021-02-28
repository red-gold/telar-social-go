package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/constants"
	domain "github.com/red-gold/ts-serverless/micros/posts/dto"
	models "github.com/red-gold/ts-serverless/micros/posts/models"
	service "github.com/red-gold/ts-serverless/micros/posts/services"
)

// CreatePostHandle handle create a new post
func CreatePostHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.CreatePostModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal CreatePostModel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelMarshalError", errorMessage)}, nil
		}

		// Create service
		postService, serviceErr := service.NewPostService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("post Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("postServiceError", errorMessage)}, nil
		}
		var newAlbum *domain.PostAlbum = nil
		if model.PostTypeId == constants.PostConstAlbum.Parse() || model.PostTypeId == constants.PostConstPhotoGallery.Parse() {
			newAlbum = &domain.PostAlbum{
				Count:   model.Album.Count,
				Cover:   model.Album.Cover,
				CoverId: model.Album.CoverId,
				Photos:  model.Album.Photos,
				Title:   model.Album.Title,
			}
		}
		newPost := &domain.Post{
			ObjectId:         model.ObjectId,
			PostTypeId:       model.PostTypeId,
			OwnerUserId:      req.UserID,
			Score:            model.Score,
			Votes:            make(map[string]bool),
			ViewCount:        model.ViewCount,
			Body:             model.Body,
			OwnerDisplayName: req.DisplayName,
			OwnerAvatar:      req.Avatar,
			Tags:             model.Tags,
			CommentCounter:   model.CommentCounter,
			Image:            model.Image,
			ImageFullPath:    model.ImageFullPath,
			Video:            model.Video,
			Thumbnail:        model.Thumbnail,
			Album:            newAlbum,
			DisableComments:  model.DisableComments,
			DisableSharing:   model.DisableSharing,
			Deleted:          model.Deleted,
			DeletedDate:      model.DeletedDate,
			CreatedDate:      utils.UTCNowUnix(),
			LastUpdated:      model.LastUpdated,
			AccessUserList:   model.AccessUserList,
			Permission:       model.Permission,
			Version:          model.Version,
		}

		if err := postService.SavePost(newPost); err != nil {
			errorMessage := fmt.Sprintf("Save new post error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("saveNewPostError", errorMessage)}, nil
		}

		return handler.Response{
			Body:       []byte(fmt.Sprintf(`{"success": true, "objectId": "%s"}`, newPost.ObjectId.String())),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// InitPostIndexHandle handle create a new post
func InitPostIndexHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create service
		postService, serviceErr := service.NewPostService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("post Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("postServiceError", errorMessage)}, nil
		}

		postIndexMap := make(map[string]interface{})
		postIndexMap["body"] = "text"
		postIndexMap["objectId"] = 1
		if err := postService.CreatePostIndex(postIndexMap); err != nil {
			errorMessage := fmt.Sprintf("Create post index Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("createPostIndexError", errorMessage)}, nil
		}

		return handler.Response{
			Body:       []byte(fmt.Sprintf(`{"success": true}`)),
			StatusCode: http.StatusOK,
		}, nil
	}
}
