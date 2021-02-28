package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	uuid "github.com/gofrs/uuid"
	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	utils "github.com/red-gold/telar-core/utils"
	"github.com/red-gold/ts-serverless/constants"
	models "github.com/red-gold/ts-serverless/micros/posts/models"
	service "github.com/red-gold/ts-serverless/micros/posts/services"
)

// QueryPostHandle handle query on post
func QueryPostHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		postService, serviceErr := service.NewPostService(db)
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
		postTypeIdParam := query.Get("type")

		var ownerUserIdList []uuid.UUID

		if ownerUserIdParam != "" {
			for _, userIdParam := range strings.Split(ownerUserIdParam, ",") {

				parsedUUID, uuidErr := uuid.FromString(userIdParam)

				if uuidErr != nil {
					return handler.Response{StatusCode: http.StatusInternalServerError}, uuidErr
				}

				ownerUserIdList = append(ownerUserIdList, parsedUUID)
			}
		}

		var postTypeId *int = nil
		if postTypeIdParam != "" {

			parsedType, strErr := strconv.Atoi(postTypeIdParam)
			if strErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, strErr
			}
			postTypeId = &parsedType
		}
		page := 0
		if pageParam != "" {
			var strErr error
			page, strErr = strconv.Atoi(pageParam)
			if strErr != nil {
				return handler.Response{StatusCode: http.StatusInternalServerError}, strErr
			}
		}
		postList, err := postService.QueryPost(searchParam, ownerUserIdList, postTypeId, "created_date", int64(page))
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		body, marshalErr := json.Marshal(postList)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling postList: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("postListMarshalError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// GetPostHandle handle get a post
func GetPostHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		postService, serviceErr := service.NewPostService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}
		postId := req.GetParamByName("postId")
		postUUID, uuidErr := uuid.FromString(postId)
		if uuidErr != nil {
			return handler.Response{StatusCode: http.StatusBadRequest,
					Body: utils.MarshalError("parseUUIDError", "Can not parse post id!")},
				nil
		}

		foundPost, err := postService.FindById(postUUID)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
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

		if foundPost.PostTypeId == constants.PostConstAlbum.Parse() || foundPost.PostTypeId == constants.PostConstPhotoGallery.Parse() {
			postModel.Album = models.PostAlbumModel{
				Count:   foundPost.Album.Count,
				Cover:   foundPost.Album.Cover,
				CoverId: foundPost.Album.CoverId,
				Photos:  foundPost.Album.Photos,
				Title:   foundPost.Album.Title,
			}
		}

		body, marshalErr := json.Marshal(postModel)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling postModel: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("marshalPostModelError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}
