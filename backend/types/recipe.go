package types

import "time"

type Recipe struct {
	Id               uint32 `json:"id"`
	SimpleStruct     bool   `json:"simple"`
	IsFork           bool   `json:"isFork"`
	OriginalRecipeId uint32 `json:"originalRecipeId"`
	IsPlaceholder    bool   `json:"isPlaceholder"`
	SourceUrl        string `json:"sourceUrl"`

	OwnerUserId    NullInt32             `json:"-"`
	LastEditUserId NullInt32             `json:"-"`
	User           NullUserProfileSimple `json:"user"`

	AiGenerated bool `json:"aiGenerated"`

	UserLocale   string                        `json:"userLocale"`
	Localization map[string]RecipeLocalization `json:"localization"`

	Categories  []RecipeCategoryitem `json:"categories"`
	Pictures    []Picture            `json:"pictures"`
	Preparation []Preparation        `json:"preparation"`

	ServingsCount uint8             `json:"servingsCount"`
	Difficulty    uint8             `json:"difficulty"`
	Statistics    RecipeStatistics  `json:"statistics"`
	Timing        PreparationTiming `json:"timing"`

	SharedInternal bool `json:"sharedInternal"`
	SharedPublic   bool `json:"sharedPublic"`

	AiTranslatedTime NullTime  `json:"localized"`
	CreatedTime      time.Time `json:"created"`
	EditedByUserTime NullTime  `json:"edited"`
	ModifiedTime     time.Time `json:"modified"`
	PublishedTime    NullTime  `json:"published"`
}

type RecipeCategoryitem struct {
	ItemId  uint16                `json:"categoryitem"`
	UserId  NullUserProfileSimple `json:"user"`
	Created time.Time             `json:"created"`
}

type RecipeLocalization struct {
	Title             string `json:"title"`
	Description       string `json:"description"`
	SourceDescription string `json:"sourceDescription"`
}

type RecipeStatistics struct {
	Steps   uint8                `json:"steps"`
	Views   uint32               `json:"views"`
	Cooked  uint32               `json:"cooked"`
	Votes   RecipeStatisticsItem `json:"votes"`
	Ratings RecipeStatisticsItem `json:"ratings"`
}

type RecipeStatisticsItem struct {
	Avg   float32 `json:"avg"`
	Count uint32  `json:"count"`
}

type Ingredient struct {
	Id           uint64                            `json:"id"`
	SortIndex    uint16                            `json:"index"`
	Quantity     NullFloat64                       `json:"quantity"`
	UnitId       NullInt32                         `json:"unitId"`
	Localization map[string]IngredientLocalization `json:"localization"`
}

type IngredientLocalization struct {
	Title string `json:"title"`
}

type Picture struct {
	Id           uint32                         `json:"id"`
	RecipeId     uint32                         `json:"-"`
	User         NullUserProfileSimple          `json:"user"`
	Index        uint8                          `json:"index"`
	Localization map[string]PictureLocalization `json:"localization"`
	FileName     string                         `json:"filename"`
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

type UserProfile struct {
	Id                    int                   `json:"id" db:"user_id"`
	UserName              string                `json:"username" db:"cloudid"`
	DisplayName           string                `json:"displayname" db:"clouddisplayname"`
	NcEnabled             bool                  `json:"-" db:"cloudenabled"`
	FirstName             string                `json:"firstname" db:"firstname"`
	LastName              string                `json:"lastname" db:"lastname"`
	Enabled               bool                  `json:"enabled" db:"enabled"`
	Admin                 bool                  `json:"admin" db:"admin"`
	Email                 NullString            `json:"email" db:"email"`
	EmailValidationPhrase NullString            `json:"-" db:"email_validationphrase"`
	EmailValidated        NullTime              `json:"-" db:"email_validated"`
	NcSyncTime            NullTime              `json:"-" db:"cloudsync"`
	NcSyncStatus          int16                 `json:"-" db:"cloudsync_status"`
	Created               time.Time             `json:"created" db:"created"`
	Modified              NullTime              `json:"-" db:"modified"`
	Groups                []Group               `json:"groups"`
	SimpleProfile         NullUserProfileSimple `json:"-"`
}

type UserProfileSimple struct {
	Id          int    `json:"id" db:"user_id"`
	DisplayName string `json:"displayname" db:"clouddisplayname"`
}

type Group struct {
	Id          int       `json:"id" db:"id"`
	DisplayName string    `json:"displayname" db:"displayname"`
	Name        string    `json:"name" db:"ncname"`
	GrantAccess bool      `json:"-" db:"access_granted"`
	GrantAdmin  string    `json:"-" db:"is_admin"`
	Created     time.Time `json:"-" db:"created"`
	Modified    NullTime  `json:"-" db:"modified"`
}
