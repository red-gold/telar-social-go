package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	utils "github.com/red-gold/telar-core/utils"
	models "github.com/red-gold/ts-serverless/src/models/gallery"
	service "github.com/red-gold/ts-serverless/src/services/gallery"
	uuid "github.com/satori/go.uuid"
)

// QueryMediaHandle handle query on media
func QueryMediaHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		mediaService, serviceErr := service.NewMediaService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		var query *url.Values
		if len(req.QueryString) > 0 {
			q, err := url.ParseQuery(string(req.QueryString))
			if err != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, err
			}
			query = &q
		}
		searchParam := query.Get("search")
		pageParam := query.Get("page")
		ownerUserIdParam := query.Get("owner")
		mediaTypeIdParam := query.Get("type")

		var ownerUserId *uuid.UUID = nil
		if ownerUserIdParam != "" {

			parsedUUID, uuidErr := uuid.FromString(ownerUserIdParam)

			if uuidErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, uuidErr
			}

			ownerUserId = &parsedUUID
		}

		var mediaTypeId *int = nil
		if mediaTypeIdParam != "" {

			parsedType, strErr := strconv.Atoi(mediaTypeIdParam)
			if strErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, strErr
			}
			mediaTypeId = &parsedType
		}
		page := 0
		if pageParam != "" {
			var strErr error
			page, strErr = strconv.Atoi(pageParam)
			if strErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, strErr
			}
		}
		mediaList, err := mediaService.QueryMedia(searchParam, ownerUserId, mediaTypeId, "created_date", int64(page))
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(mediaList)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling mediaList: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("mediaListMarshalError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// QueryAlbumHandle handle query on media
func QueryAlbumHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		mediaService, serviceErr := service.NewMediaService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		var query *url.Values
		if len(req.QueryString) > 0 {
			q, err := url.ParseQuery(string(req.QueryString))
			if err != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, err
			}
			query = &q
		}
		pageParam := query.Get("page")
		limitParam := query.Get("limit")
		albumIdParam := query.Get("album")

		var albumId *uuid.UUID = nil
		if albumIdParam != "" {

			parsedUUID, uuidErr := uuid.FromString(albumIdParam)

			if uuidErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, uuidErr
			}

			albumId = &parsedUUID
		}

		page := 0
		if pageParam != "" {
			var strErr error
			page, strErr = strconv.Atoi(pageParam)
			if strErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, strErr
			}
		}

		limit := 0
		if limitParam != "" {
			var strErr error
			limit, strErr = strconv.Atoi(limitParam)
			if strErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, strErr
			}
		}
		mediaList, err := mediaService.QueryAlbum(req.UserID, albumId, int64(page), int64(limit), "created_date")
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(mediaList)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling mediaList: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("mediaListMarshalError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// GetMediaHandle handle get a media
func GetMediaHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// params from /medias/id/:mediaId
		mediaId := req.GetParamByName("mediaId")
		mediaUUID, uuidErr := uuid.FromString(mediaId)
		if uuidErr != nil {
			return handler.Response{StatusCode: http.StatusBadRequest,
					Body: utils.MarshalError("parseUUIDError", "Can not parse media id!")},
				nil
		}

		// Create service
		mediaService, serviceErr := service.NewMediaService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		foundMedia, err := mediaService.FindById(mediaUUID)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		mediaModel := models.MediaModel{
			ObjectId:       foundMedia.ObjectId,
			DeletedDate:    foundMedia.DeletedDate,
			CreatedDate:    foundMedia.CreatedDate,
			Thumbnail:      foundMedia.Thumbnail,
			URL:            foundMedia.URL,
			FullPath:       foundMedia.FullPath,
			Caption:        foundMedia.Caption,
			FileName:       foundMedia.FileName,
			Directory:      foundMedia.Directory,
			OwnerUserId:    foundMedia.OwnerUserId,
			LastUpdated:    foundMedia.LastUpdated,
			AlbumId:        foundMedia.AlbumId,
			Width:          foundMedia.Width,
			Height:         foundMedia.Height,
			Meta:           foundMedia.Meta,
			AccessUserList: foundMedia.AccessUserList,
			Permission:     foundMedia.Permission,
			Deleted:        foundMedia.Deleted,
		}

		body, marshalErr := json.Marshal(mediaModel)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling mediaModel: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("marshalMediaModelError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// GetMediaByDirectoryHandle handle get media list by directory
func GetMediaByDirectoryHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {

		// params from /medias/dir/:dir
		dirName := req.GetParamByName("dir")
		if dirName == "" {
			errorMessage := fmt.Sprintf("Directory name is required!")
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("dirNameRequired", errorMessage)}, nil
		}

		// Create service
		mediaService, serviceErr := service.NewMediaService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		foundMediaList, err := mediaService.FindByDirectory(req.UserID, dirName, 0, 0)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(foundMediaList)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling mediaModel: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("marshalMediaModelError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}
