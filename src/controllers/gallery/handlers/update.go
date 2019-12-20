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

// UpdateMediaHandle handle create a new media
func UpdateMediaHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// Create the model object
		var model models.MediaModel
		if err := json.Unmarshal(req.Body, &model); err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		// Create service
		mediaService, serviceErr := service.NewMediaService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		updatedMedia := &domain.Media{
			ObjectId:       model.ObjectId,
			DeletedDate:    0,
			CreatedDate:    model.CreatedDate,
			Thumbnail:      model.Thumbnail,
			URL:            model.URL,
			FullPath:       model.FullPath,
			Caption:        model.Caption,
			FileName:       model.FileName,
			Directory:      model.Directory,
			OwnerUserId:    req.UserID,
			LastUpdated:    utils.UTCNowUnix(),
			AlbumId:        model.AlbumId,
			Width:          model.Width,
			Height:         model.Height,
			Meta:           model.Meta,
			AccessUserList: model.AccessUserList,
			Permission:     model.Permission,
			Deleted:        model.Deleted,
		}

		if err := mediaService.UpdateMediaById(updatedMedia); err != nil {
			errorMessage := fmt.Sprintf("Update Media Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("updateMediaError", errorMessage)}, nil
		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
