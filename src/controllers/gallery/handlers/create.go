package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	domain "github.com/red-gold/ts-serverless/src/domain/gallery"
	models "github.com/red-gold/ts-serverless/src/models/gallery"
	service "github.com/red-gold/ts-serverless/src/services/gallery"
)

// CreateMediaHandle handle create a new media
func CreateMediaHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.CreateMediaModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal CreateMediaModel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelMarshalError", errorMessage)}, nil
		}

		// Create service
		mediaService, serviceErr := service.NewMediaService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("media Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("mediaServiceError", errorMessage)}, nil
		}

		newMedia := &domain.Media{
			ObjectId:       model.ObjectId,
			DeletedDate:    0,
			CreatedDate:    utils.UTCNowUnix(),
			Thumbnail:      model.Thumbnail,
			URL:            model.URL,
			FullPath:       model.FullPath,
			Caption:        model.Caption,
			FileName:       model.FileName,
			Directory:      model.Directory,
			OwnerUserId:    req.UserID,
			LastUpdated:    0,
			AlbumId:        model.AlbumId,
			Width:          model.Width,
			Height:         model.Height,
			Meta:           model.Meta,
			AccessUserList: model.AccessUserList,
			Permission:     model.Permission,
			Deleted:        false,
		}

		if err := mediaService.SaveMedia(newMedia); err != nil {
			errorMessage := fmt.Sprintf("Save Media Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("saveMediaError", errorMessage)}, nil
		}

		return handler.Response{
			Body:       []byte(fmt.Sprintf(`{"success": true, "objectId": "%s"}`, newMedia.ObjectId.String())),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// CreateMediaListHandle handle create a new media
func CreateMediaListHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.CreateMediaListModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			errorMessage := fmt.Sprintf("Unmarshal CreateMediaModel Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("modelMarshalError", errorMessage)}, nil
		}

		// Create service
		mediaService, serviceErr := service.NewMediaService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("media Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("mediaServiceError", errorMessage)}, nil
		}
		var mediaList []domain.Media
		for _, media := range model.List {

			newMedia := domain.Media{
				ObjectId:       media.ObjectId,
				DeletedDate:    0,
				CreatedDate:    utils.UTCNowUnix(),
				Thumbnail:      media.Thumbnail,
				URL:            media.URL,
				FullPath:       media.FullPath,
				Caption:        media.Caption,
				FileName:       media.FileName,
				Directory:      media.Directory,
				OwnerUserId:    req.UserID,
				LastUpdated:    0,
				AlbumId:        media.AlbumId,
				Width:          media.Width,
				Height:         media.Height,
				Meta:           media.Meta,
				AccessUserList: media.AccessUserList,
				Permission:     media.Permission,
				Deleted:        false,
			}
			mediaList = append(mediaList, newMedia)
		}

		if err := mediaService.SaveManyMedia(mediaList); err != nil {
			errorMessage := fmt.Sprintf("Save Media Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("saveMediaError", errorMessage)}, nil
		}

		return handler.Response{
			Body:       []byte(fmt.Sprintf(`{"success": true}`)),
			StatusCode: http.StatusOK,
		}, nil
	}
}
