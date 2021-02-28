package constants

type PostConst int

const (
	PostConstText PostConst = iota
	PostConstPhoto
	PostConstVideo
	PostConstPhotoGallery
	PostConstAlbum
)

func (pc PostConst) Parse() int {
	return int(pc)
}
