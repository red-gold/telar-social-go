package handlers

import (
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-core/utils"
	service "github.com/red-gold/ts-serverless/micros/gallery/services"
	uuid "github.com/satori/go.uuid"
)

// DeleteMediaHandle handle delete a media
func DeleteMediaHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /medias/id/:mediaId
		mediaId := req.GetParamByName("mediaId")
		if mediaId == "" {
			errorMessage := fmt.Sprintf("Media Id is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("mediaIdRequired", errorMessage)}, nil
		}
		fmt.Printf("\n Media ID: %s", mediaId)
		mediaUUID, uuidErr := uuid.FromString(mediaId)
		if uuidErr != nil {
			errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("uuidError", errorMessage)}, nil
		}
		fmt.Printf("\n Media UUID: %s", mediaUUID)
		// Create service
		mediaService, serviceErr := service.NewMediaService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("Media Service Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("mediaServiceError", errorMessage)}, nil

		}

		if err := mediaService.DeleteMediaByOwner(req.UserID, mediaUUID); err != nil {
			errorMessage := fmt.Sprintf("Delete Media Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("deleteMediaError", errorMessage)}, nil

		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// DeleteDirectoryHandle handle delete a media
func DeleteDirectoryHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /medias/dir/:dir
		dirName := req.GetParamByName("dir")
		if dirName == "" {
			errorMessage := fmt.Sprintf("Directory name is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("dirNameRequired", errorMessage)}, nil
		}
		fmt.Printf("\n Directory ID: %s", dirName)

		// Create service
		mediaService, serviceErr := service.NewMediaService(db)
		if serviceErr != nil {
			errorMessage := fmt.Sprintf("Media Service Error %s", serviceErr.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("mediaServiceError", errorMessage)}, nil

		}

		if err := mediaService.DeleteMediaByDirectory(req.UserID, dirName); err != nil {
			errorMessage := fmt.Sprintf("Delete Media Error %s", err.Error())
			return handler.Response{StatusCode: http.StatusInternalServerError, Body: utils.MarshalError("deleteMediaError", errorMessage)}, nil

		}
		return handler.Response{
			Body:       []byte(`{"success": true}`),
			StatusCode: http.StatusOK,
		}, nil
	}
}
