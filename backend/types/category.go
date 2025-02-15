package types

import "time"

// Category represents a category in the database
type Category struct {
	Id       uint16         `json:"id"`
	Name     string         `json:"name"`
	Icon     string         `json:"icon"`
	Modified time.Time      `json:"modified"`
	Items    []CategoryItem `json:"items"`
}

type CategoryItem struct {
	Id       uint16    `json:"id"`
	Name     string    `json:"name"`
	Icon     string    `json:"icon"`
	Modified time.Time `json:"modified"`
}
