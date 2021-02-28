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
	uuid "github.com/satori/go.uuid"
)

// UpdatePostHandle handle create a new post
func UpdatePostHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.PostModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		// Create service
		postService, serviceErr := service.NewPostService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}
		var updatedAlbum *domain.PostAlbum = nil
		if model.PostTypeId == constants.PostConstAlbum.Parse() || model.PostTypeId == constants.PostConstPhotoGallery.Parse() {
			updatedAlbum = &domain.PostAlbum{
				Count:   model.Album.Count,
				Cover:   model.Album.Cover,
				CoverId: model.Album.CoverId,
				Photos:  model.Album.Photos,
				Title:   model.Album.Title,
			}
		}
		updatedPost := &domain.Post{
			ObjectId:         model.ObjectId,
			PostTypeId:       model.PostTypeId,
			OwnerUserId:      req.UserID,
			Score:            model.Score,
			Votes:            model.Votes,
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
			Album:            updatedAlbum,
			DisableComments:  model.DisableComments,
			DisableSharing:   model.DisableSharing,
			Deleted:          model.Deleted,
			DeletedDate:      model.DeletedDate,
			CreatedDate:      model.CreatedDate,
			LastUpdated:      utils.UTCNowUnix(),
			AccessUserList:   model.AccessUserList,
			Permission:       model.Permission,
			Version:          model.Version,
		}

		if err := postService.UpdatePostById(updatedPost); err != nil {
			errorMessage := fmt.Sprintf("Update Post Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("updatePostError", errorMessage)}, nil
		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// IncrementScoreHandle handle create a new post
func IncrementScoreHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /posts/score/+1/:postId
		postId := req.GetParamByName("postId")
		if postId == "" {
			errorMessage := fmt.Sprintf("Post Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("postIdRequired", errorMessage)}, nil
		}

		fmt.Printf("\n Post ID: %s", postId)
		postUUID, uuidErr := uuid.FromString(postId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}
		fmt.Printf("\n Post UUID: %s", postUUID)
		// Create service
		postService, serviceErr := service.NewPostService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		fmt.Printf("\nReqUserID: %s\n", req.UserID)
		postService.IncrementScoreCount(postUUID, req.UserID)
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// DecrementScoreHandle handle create a new post
func DecrementScoreHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /posts/score/-1/:postId
		postId := req.GetParamByName("postId")
		if postId == "" {
			errorMessage := fmt.Sprintf("Post Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("postIdRequired", errorMessage)}, nil
		}

		fmt.Printf("\n Post ID: %s", postId)
		postUUID, uuidErr := uuid.FromString(postId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}
		fmt.Printf("\n Post UUID: %s", postUUID)
		// Create service
		postService, serviceErr := service.NewPostService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		postService.DecrementScoreCount(postUUID, req.UserID)
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// IncrementCommentHandle handle create a new post
func IncrementCommentHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /comment/+1/:postId
		postId := req.GetParamByName("postId")
		if postId == "" {
			errorMessage := fmt.Sprintf("Post Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("postIdRequired", errorMessage)}, nil
		}

		fmt.Printf("\n Post ID: %s", postId)
		postUUID, uuidErr := uuid.FromString(postId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}
		fmt.Printf("\n Post UUID: %s", postUUID)
		// Create service
		postService, serviceErr := service.NewPostService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		postService.IncrementCommentCount(postUUID)
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// DecrementCommentHandle handle create a new post
func DecrementCommentHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /comment/-1/:postId
		postId := req.GetParamByName("postId")
		if postId == "" {
			errorMessage := fmt.Sprintf("Post Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("postIdRequired", errorMessage)}, nil
		}

		fmt.Printf("\n Post ID: %s", postId)
		postUUID, uuidErr := uuid.FromString(postId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}
		fmt.Printf("\n Post UUID: %s", postUUID)
		// Create service
		postService, serviceErr := service.NewPostService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		postService.DecerementCommentCount(postUUID)
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// UpdatePostProfileHandle handle create a new post
func UpdatePostProfileHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create service
		postService, serviceErr := service.NewPostService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		postService.UpdatePostProfile(req.UserID, req.DisplayName, req.Avatar)
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
