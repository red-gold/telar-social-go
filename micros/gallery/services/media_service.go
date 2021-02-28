package service

import (
	"fmt"

	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/config"
	coreData "github.com/red-gold/telar-core/data"
	repo "github.com/red-gold/telar-core/data"
	"github.com/red-gold/telar-core/data/mongodb"
	mongoRepo "github.com/red-gold/telar-core/data/mongodb"
	"github.com/red-gold/telar-core/utils"
	dto "github.com/red-gold/ts-serverless/micros/gallery/dto"
)

// MediaService handlers with injected dependencies
type MediaServiceImpl struct {
	MediaRepo repo.Repository
}

// NewMediaService initializes MediaService's dependencies and create new MediaService struct
func NewMediaService(db interface{}) (MediaService, error) {

	mediaService := &MediaServiceImpl{}

	switch *config.AppConfig.DBType {
	case config.DB_MONGO:

		mongodb := db.(mongodb.MongoDatabase)
		mediaService.MediaRepo = mongoRepo.NewDataRepositoryMongo(mongodb)

	}

	return mediaService, nil
}

// SaveMedia save the media
func (s MediaServiceImpl) SaveMedia(media *dto.Media) error {

	if media.ObjectId == uuid.Nil {
		var uuidErr error
		media.ObjectId, uuidErr = uuid.NewV4()
		if uuidErr != nil {
			return uuidErr
		}
	}

	if media.CreatedDate == 0 {
		media.CreatedDate = utils.UTCNowUnix()
	}

	result := <-s.MediaRepo.Save(mediaCollectionName, media)

	return result.Error
}

// SaveManyMedia save the media
func (s MediaServiceImpl) SaveManyMedia(medias []dto.Media) error {

	// https://github.com/golang/go/wiki/InterfaceSlice
	var interfaceSlice []interface{} = make([]interface{}, len(medias))
	for i, d := range medias {
		if d.ObjectId == uuid.Nil {
			var uuidErr error
			d.ObjectId, uuidErr = uuid.NewV4()
			if uuidErr != nil {
				return uuidErr
			}
		}

		if d.CreatedDate == 0 {
			d.CreatedDate = utils.UTCNowUnix()
		}
		interfaceSlice[i] = d
	}
	result := <-s.MediaRepo.SaveMany(mediaCollectionName, interfaceSlice)

	return result.Error
}

// FindOneMedia get one media
func (s MediaServiceImpl) FindOneMedia(filter interface{}) (*dto.Media, error) {

	result := <-s.MediaRepo.FindOne(mediaCollectionName, filter)
	if result.Error() != nil {
		return nil, result.Error()
	}

	var mediaResult dto.Media
	errDecode := result.Decode(&mediaResult)
	if errDecode != nil {
		return nil, fmt.Errorf("Error docoding on dto.Media")
	}
	return &mediaResult, nil
}

// FindMediaList get all medias by filter
func (s MediaServiceImpl) FindMediaList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Media, error) {

	result := <-s.MediaRepo.Find(mediaCollectionName, filter, limit, skip, sort)
	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var mediaList []dto.Media
	for result.Next() {
		var media dto.Media
		errDecode := result.Decode(&media)
		if errDecode != nil {
			return nil, fmt.Errorf("Error docoding on dto.Media")
		}
		mediaList = append(mediaList, media)
	}

	return mediaList, nil
}

// QueryMedia get all medias by query
func (s MediaServiceImpl) QueryMedia(search string, ownerUserId *uuid.UUID, mediaTypeId *int, sortBy string, page int64) ([]dto.Media, error) {
	sortMap := make(map[string]int)
	sortMap[sortBy] = -1
	skip := numberOfItems * (page - 1)
	limit := numberOfItems

	filter := make(map[string]interface{})
	if search != "" {
		filter["$text"] = coreData.SearchOperator{Search: search}
	}
	if ownerUserId != nil {
		filter["ownerUserId"] = *ownerUserId
	}
	if mediaTypeId != nil {
		filter["mediaTypeId"] = *mediaTypeId
	}
	fmt.Println(filter)
	result, err := s.FindMediaList(filter, limit, skip, sortMap)

	return result, err
}

// FindByOwnerUserId find by owner user id
func (s MediaServiceImpl) FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Media, error) {
	sortMap := make(map[string]int)
	sortMap["created_date"] = -1
	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		OwnerUserId: ownerUserId,
	}
	return s.FindMediaList(filter, 0, 0, sortMap)
}

// FindById find by media id
func (s MediaServiceImpl) FindById(objectId uuid.UUID) (*dto.Media, error) {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}
	return s.FindOneMedia(filter)
}

// UpdateMedia update the media
func (s MediaServiceImpl) UpdateMedia(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error {

	result := <-s.MediaRepo.Update(mediaCollectionName, filter, data, opts...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateMedia update the media
func (s MediaServiceImpl) UpdateMediaById(data *dto.Media) error {
	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    data.ObjectId,
		OwnerUserId: data.OwnerUserId,
	}

	updateOperator := coreData.UpdateOperator{
		Set: data,
	}
	err := s.UpdateMedia(filter, updateOperator)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMedia delete media by filter
func (s MediaServiceImpl) DeleteMedia(filter interface{}) error {

	result := <-s.MediaRepo.Delete(mediaCollectionName, filter, true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteMedia delete media by ownerUserId and mediaId
func (s MediaServiceImpl) DeleteMediaByOwner(ownerUserId uuid.UUID, mediaId uuid.UUID) error {

	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    mediaId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteMedia(filter)
	if err != nil {
		return err
	}
	return nil
}

// DeleteManyMedia delete many media by filter
func (s MediaServiceImpl) DeleteManyMedia(filter interface{}) error {

	result := <-s.MediaRepo.Delete(mediaCollectionName, filter, false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateMediaIndex create index for media search.
func (s MediaServiceImpl) CreateMediaIndex(indexes map[string]interface{}) error {
	result := <-s.MediaRepo.CreateIndex(mediaCollectionName, indexes)
	return result
}

// FindByDirectory find by directory
func (s MediaServiceImpl) FindByDirectory(ownerUserId uuid.UUID, directory string, limit int64, skip int64) ([]dto.Media, error) {
	sortMap := make(map[string]int)
	sortMap["created_date"] = -1
	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
		Directory   string    `json:"directory" bson:"directory"`
	}{
		OwnerUserId: ownerUserId,
		Directory:   directory,
	}
	return s.FindMediaList(filter, limit, skip, sortMap)
}

// QueryAlbum query media by albumId
func (s MediaServiceImpl) QueryAlbum(ownerUserId uuid.UUID, albumId *uuid.UUID, page int64, limit int64, sortBy string) ([]dto.Media, error) {
	sortMap := make(map[string]int)
	sortMap[sortBy] = -1
	skip := numberOfItems * (page - 1)
	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
		AlbumId     uuid.UUID `json:"albumId" bson:"albumId"`
	}{
		OwnerUserId: ownerUserId,
		AlbumId:     *albumId,
	}
	return s.FindMediaList(filter, limit, skip, sortMap)
}

// DeleteMediaByDirectory delete media by ownerUserId and mediaId
func (s MediaServiceImpl) DeleteMediaByDirectory(ownerUserId uuid.UUID, directory string) error {

	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
		Directory   string    `json:"directory" bson:"directory"`
	}{
		OwnerUserId: ownerUserId,
		Directory:   directory,
	}
	err := s.DeleteManyMedia(filter)
	if err != nil {
		return err
	}
	return nil
}
