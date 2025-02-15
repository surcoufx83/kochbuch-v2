package types

import "time"

type Recipe struct {
	Id               uint32 `json:"id"`
	IsFork           bool   `json:"isFork"`
	OriginalRecipeId uint32 `json:"originalRecipeId"`
	IsPlaceholder    bool   `json:"isPlaceholder"`
	SourceUrl        string `json:"sourceUrl"`

	OwnerUserId    uint16 `json:"-"`
	LastEditUserId uint16 `json:"-"`

	AiGenerated bool `json:"aiGenerated"`
	AiLocalized bool `json:"aiLocalized"`

	Localization map[string]RecipeLocalization `json:"localization"`

	Categories []uint16 `json:"categories"`

	ServingsCount uint8 `json:"servingsCount"`
	Difficulty    uint8 `json:"difficulty"`

	SharedInternal bool `json:"sharedInternal"`
	SharedPublic   bool `json:"sharedPublic"`

	CreatedTime   time.Time `json:"created"`
	ModifiedTime  time.Time `json:"modified"`
	PublishedTime time.Time `json:"published"`
}

type RecipeLocalization struct {
	UserGeneratedtime bool   `json:"userGenerated"`
	Default           bool   `json:"default"`
	Language          string `json:"language"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	SourceDescription string `json:"sourceDescription"`
}

type Ingredient struct {
	Id           uint64                            `json:"id"`
	StepId       uint64                            `json:"stepId"`
	SortIndex    uint8                             `json:"sortIndex"`
	Quantity     float32                           `json:"quantity"`
	Unit         uint8                             `json:"unit"`
	Localization map[string]IngredientLocalization `json:"localization"`
}

type IngredientLocalization struct {
	UserGeneratedtime bool   `json:"userGenerated"`
	Default           bool   `json:"default"`
	Language          string `json:"language"`
	Title             string `json:"title"`
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
