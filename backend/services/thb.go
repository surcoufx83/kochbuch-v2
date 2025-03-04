package services

import (
	"kochbuch-v2-backend/types"
	"time"
)

var (
	ThumbnailGenerationRequests []PictureRequiresThumbnail
)

type PictureRequiresThumbnail struct {
	RecipeId  uint32
	PictureId uint32
	Index     uint8
	Picture   *types.Picture
}

func ThbAutoGenerator() {
	for {
		if len(ThumbnailGenerationRequests) == 0 {
			// log.Println("No picture require thumbnails")
			time.Sleep(5 * time.Minute)
		} else {
			resizePictures()
			time.Sleep(5 * time.Second)
		}
	}
}

func resizePictures() {
	if len(ThumbnailGenerationRequests) == 0 {
		return
	}

	var req PictureRequiresThumbnail
	req, ThumbnailGenerationRequests = ThumbnailGenerationRequests[0], ThumbnailGenerationRequests[1:]
	GenerateResizedPictureVersions(req.RecipeId, req.PictureId)

}
