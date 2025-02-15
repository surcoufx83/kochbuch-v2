package types

import "time"

// Category represents a category in the database
type Category struct {
	Id           uint16                      `json:"id"`
	Localization map[string]NameLocalization `json:"localization"`
	Icon         string                      `json:"icon"`
	Modified     time.Time                   `json:"modified"`
	Items        map[uint16]CategoryItem     `json:"items"`
}

type CategoryItem struct {
	Id           uint16                      `json:"id"`
	Localization map[string]NameLocalization `json:"localization"`
	Icon         string                      `json:"icon"`
	Modified     time.Time                   `json:"modified"`
}
