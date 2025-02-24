package types

import "time"

type Recipe struct {
	Id               uint32 `json:"id"`
	SimpleStruct     bool   `json:"simple"`
	IsFork           bool   `json:"isFork"`
	OriginalRecipeId uint32 `json:"originalRecipeId"`
	IsPlaceholder    bool   `json:"isPlaceholder"`
	SourceUrl        string `json:"sourceUrl"`

	OwnerUserId    NullInt32 `json:"-"`
	LastEditUserId NullInt32 `json:"-"`

	AiGenerated bool `json:"aiGenerated"`
	AiLocalized bool `json:"aiLocalized"`

	UserLocale   string                        `json:"userLocale"`
	Localization map[string]RecipeLocalization `json:"localization"`

	Categories  []uint16      `json:"categories"`
	Pictures    []Picture     `json:"pictures"`
	Preparation []Preparation `json:"preparation"`

	ServingsCount uint8 `json:"servingsCount"`
	Difficulty    uint8 `json:"difficulty"`

	SharedInternal bool `json:"sharedInternal"`
	SharedPublic   bool `json:"sharedPublic"`

	CreatedTime   time.Time `json:"created"`
	ModifiedTime  time.Time `json:"modified"`
	PublishedTime NullTime  `json:"published"`
}

type RecipeLocalization struct {
	Title             string `json:"title"`
	Description       string `json:"description"`
	SourceDescription string `json:"sourceDescription"`
}

type Ingredient struct {
	Id            uint64                            `json:"id"`
	PreparationId uint64                            `json:"-"`
	RecipeId      uint32                            `json:"-"`
	SortIndex     uint8                             `json:"index"`
	Quantity      float32                           `json:"quantity"`
	UnitId        uint8                             `json:"unitId"`
	Localization  map[string]IngredientLocalization `json:"localization"`
}

type IngredientLocalization struct {
	Title string `json:"title"`
}

type Picture struct {
	Id           uint32                         `json:"id"`
	RecipeId     uint32                         `json:"-"`
	UserId       NullInt32                      `json:"userId"`
	Index        uint8                          `json:"index"`
	Localization map[string]PictureLocalization `json:"localization"`
	FileName     string                         `json:"-"`
	FullPath     string                         `json:"-"`
	Uploaded     NullTime                       `json:"uploaded"`
	Dimension    PictureDimension               `json:"size"`
}

type PictureDimension struct {
	Height int32 `json:"height"`
	Width  int32 `json:"width"`
}

type PictureLocalization struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Preparation struct {
	Id           uint64                             `json:"id"`
	RecipeId     uint32                             `json:"-"`
	Index        uint8                              `json:"index"`
	Ingredients  []Ingredient                       `json:"ingredients"`
	Localization map[string]PreparationLocalization `json:"localization"`
	Timing       PreparationTiming                  `json:"timing"`
}

type PreparationLocalization struct {
	Title        string `json:"title"`
	Instructions string `json:"instruct"`
}

type PreparationTiming struct {
	Preparing NullInt32 `json:"preparing"`
	Cooking   NullInt32 `json:"cooking"`
	Waiting   NullInt32 `json:"waiting"`
}

type Unit struct {
	Id             uint8                       `json:"id"`
	Localization   map[string]UnitLocalization `json:"localization"`
	CreatedTime    time.Time                   `json:"created"`
	ModifiedTime   time.Time                   `json:"modified"`
	ReplacedById   uint8                       `json:"replacedBy"`
	SavedAsId      uint8                       `json:"savedAs"`
	SavedAsFactor  float32                     `json:"savedAsFactor"`
	DecimalPlaces  uint8                       `json:"decimalPlaces"`
	ShowAsFraction bool                        `json:"showAsFraction"`
}

type UnitLocalization struct {
	NameSingular string `json:"singular"`
	NamePlural   string `json:"plural"`
}
