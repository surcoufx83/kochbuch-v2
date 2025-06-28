package types

import "time"

type Recipe struct {
	AiGenerated      bool                          `json:"aiGenerated"`
	AiTranslatedTime NullTime                      `json:"localized"`
	Categories       []*RecipeCategoryitem         `json:"categories"`
	CreatedTime      time.Time                     `json:"created"`
	Difficulty       uint8                         `json:"difficulty"`
	EditedByUserTime NullTime                      `json:"edited"`
	Id               uint32                        `json:"id"`
	IsFork           bool                          `json:"isFork"`
	IsPlaceholder    bool                          `json:"isPlaceholder"`
	LastEditUserId   NullInt32                     `json:"-"`
	Localization     map[string]RecipeLocalization `json:"localization"`
	ModifiedTime     time.Time                     `json:"modified"`
	OriginalRecipeId uint32                        `json:"originalRecipeId"`
	OwnerUserId      NullInt32                     `json:"-"`
	Pictures         []*Picture                    `json:"pictures"`
	Preparation      []*Preparation                `json:"preparation"`
	PublishedTime    NullTime                      `json:"published"`
	ServingsCount    uint8                         `json:"servingsCount"`
	SharedInternal   bool                          `json:"sharedInternal"`
	SharedPublic     bool                          `json:"sharedPublic"`
	SimpleStruct     bool                          `json:"simple"`
	SourceUrl        string                        `json:"sourceUrl"`
	Statistics       RecipeStatistics              `json:"statistics"`
	Timing           PreparationTiming             `json:"timing"`
	User             NullUserProfileSimple         `json:"user"`
	UserLocale       string                        `json:"userLocale"`

	SimpleRecipe RecipeSimple `json:"-"`
}

type RecipeSimple struct {
	AiGenerated      *bool                          `json:"aiGenerated"`
	AiTranslatedTime *NullTime                      `json:"localized"`
	Categories       *[]*RecipeCategoryitem         `json:"categories"`
	CreatedTime      *time.Time                     `json:"created"`
	Difficulty       *uint8                         `json:"difficulty"`
	EditedByUserTime *NullTime                      `json:"edited"`
	Id               *uint32                        `json:"id"`
	Localization     *map[string]RecipeLocalization `json:"localization"`
	ModifiedTime     *time.Time                     `json:"modified"`
	Pictures         []*Picture                     `json:"pictures"`
	PublishedTime    *NullTime                      `json:"published"`
	ServingsCount    *uint8                         `json:"servingsCount"`
	SimpleStruct     bool                           `json:"simple"`
	Statistics       *RecipeStatistics              `json:"statistics"`
	Timing           *PreparationTiming             `json:"timing"`
	User             *NullUserProfileSimple         `json:"user"`
	UserLocale       *string                        `json:"userLocale"`
}

type RecipeCategoryitem struct {
	Created time.Time             `json:"created"`
	ItemId  uint16                `json:"categoryitem"`
	UserId  NullUserProfileSimple `json:"user"`
}

type RecipeLocalization struct {
	Description       string `json:"description"`
	SourceDescription string `json:"sourceDescription"`
	Title             string `json:"title"`
}

type RecipeStatistics struct {
	Cooked  uint32               `json:"cooked"`
	Ratings RecipeStatisticsItem `json:"ratings"`
	Steps   uint8                `json:"steps"`
	Views   uint32               `json:"views"`
	Votes   RecipeStatisticsItem `json:"votes"`
}

type RecipeStatisticsItem struct {
	Avg   float32 `json:"avg"`
	Count uint32  `json:"count"`
}

type Ingredient struct {
	Id           uint64                            `json:"id"`
	Localization map[string]IngredientLocalization `json:"localization"`
	Quantity     NullFloat64                       `json:"quantity"`
	SortIndex    uint16                            `json:"index"`
	UnitId       NullInt32                         `json:"unitId"`
}

type IngredientLocalization struct {
	Title string `json:"title"`
}

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

type Preparation struct {
	Id           uint64                             `json:"id"`
	Index        uint8                              `json:"index"`
	Ingredients  []Ingredient                       `json:"ingredients"`
	Localization map[string]PreparationLocalization `json:"localization"`
	Timing       PreparationTiming                  `json:"timing"`
}

type PreparationLocalization struct {
	Instructions string `json:"instruct"`
	Title        string `json:"title"`
}

type PreparationTiming struct {
	Cooking   NullInt32 `json:"cooking"`
	Preparing NullInt32 `json:"preparing"`
	Total     NullInt32 `json:"total"`
	Waiting   NullInt32 `json:"waiting"`
}

type Unit struct {
	CreatedTime    time.Time                   `json:"created"`
	DecimalPlaces  uint8                       `json:"decimalPlaces"`
	Id             uint8                       `json:"id"`
	Localization   map[string]UnitLocalization `json:"localization"`
	ModifiedTime   time.Time                   `json:"modified"`
	ReplacedById   uint8                       `json:"replacedBy"`
	SavedAsFactor  float32                     `json:"savedAsFactor"`
	SavedAsId      uint8                       `json:"savedAs"`
	ShowAsFraction bool                        `json:"showAsFraction"`
}

type UnitLocalization struct {
	NameSingular string `json:"singular"`
	NamePlural   string `json:"plural"`
}

type UserProfile struct {
	Admin                 bool                  `json:"admin" db:"admin"`
	Created               time.Time             `json:"created" db:"created"`
	DisplayName           string                `json:"displayname" db:"clouddisplayname"`
	Email                 NullString            `json:"email" db:"email"`
	EmailValidated        NullTime              `json:"-" db:"email_validated"`
	EmailValidationPhrase NullString            `json:"-" db:"email_validationphrase"`
	Enabled               bool                  `json:"enabled" db:"enabled"`
	FirstName             string                `json:"firstname" db:"firstname"`
	Groups                []Group               `json:"groups"`
	Id                    int                   `json:"id" db:"user_id"`
	LastName              string                `json:"lastname" db:"lastname"`
	Modified              NullTime              `json:"-" db:"modified"`
	NcEnabled             bool                  `json:"-" db:"cloudenabled"`
	NcSyncStatus          int16                 `json:"-" db:"cloudsync_status"`
	NcSyncTime            NullTime              `json:"-" db:"cloudsync"`
	UserName              string                `json:"username" db:"cloudid"`
	SimpleProfile         NullUserProfileSimple `json:"-"`
}

type UserProfileSimple struct {
	Id          int    `json:"id" db:"user_id"`
	DisplayName string `json:"displayname" db:"clouddisplayname"`
}

type Group struct {
	Created     time.Time `json:"-" db:"created"`
	DisplayName string    `json:"displayname" db:"displayname"`
	GrantAccess bool      `json:"-" db:"access_granted"`
	GrantAdmin  string    `json:"-" db:"is_admin"`
	Id          int       `json:"id" db:"id"`
	Modified    NullTime  `json:"-" db:"modified"`
	Name        string    `json:"name" db:"ncname"`
}
