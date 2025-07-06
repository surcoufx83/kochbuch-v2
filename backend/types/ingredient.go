package types

type Ingredient struct {
	Id           uint64                            `json:"id"`
	Localization map[string]IngredientLocalization `json:"localization"`
	Quantity     NullFloat64                       `json:"quantity"`
	RecipeId     uint32                            `json:"-"`
	SortIndex    uint16                            `json:"index"`
	StepId       uint64                            `json:"-"`
	UnitId       NullInt32                         `json:"unitId"`
}

type IngredientLocalization struct {
	Title string `json:"title"`
}
