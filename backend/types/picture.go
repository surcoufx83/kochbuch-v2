package types

import "time"

type Picture struct {
	BaseName     string                         `json:"-"`
	Dimension    PictureDimension               `json:"size"`
	Ext          string                         `json:"-"`
	FileName     string                         `json:"filename"`
	FullPath     string                         `json:"-"`
	Id           uint32                         `json:"id"`
	Index        uint8                          `json:"index"`
	Localization map[string]PictureLocalization `json:"localization"`
	RecipeId     uint32                         `json:"-"`
	Uploaded     time.Time                      `json:"uploaded"`
	User         NullUserProfileSimple          `json:"user"`
}

type PictureDimension struct {
	Generated      NullTime `json:"thbGenerated"`
	GeneratedSizes []int    `json:"thbSizes"`
	Height         int      `json:"height"`
	Width          int      `json:"width"`
}

type PictureLocalization struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}
