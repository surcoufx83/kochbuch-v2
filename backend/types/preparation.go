package types

type Preparation struct {
	Id           uint64                             `json:"id"`
	Index        uint8                              `json:"index"`
	Ingredients  []*Ingredient                      `json:"ingredients"`
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
