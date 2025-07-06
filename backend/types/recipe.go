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
