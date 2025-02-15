package types

import "time"

// Category represents a category in the database
type Category struct {
	CategoryId   uint16         `json:"id"`
	CategoryName string         `json:"name"`
	CategoryIcon string         `json:"icon"`
	Modified     time.Time      `json:"modified"`
	Items        []CategoryItem `json:"items"`
}

type CategoryItem struct {
	CategoryItemId   uint16    `json:"id"`
	CategoryItemName string    `json:"name"`
	CategoryItemIcon string    `json:"icon"`
	Modified         time.Time `json:"modified"`
}
