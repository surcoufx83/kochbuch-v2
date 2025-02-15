package types

import "time"

// Category represents a category in the database
type Category struct {
	CategoryId   uint16         `db:"catid" json:"id"`
	CategoryName string         `db:"catname" json:"name"`
	CategoryIcon string         `db:"caticon" json:"icon"`
	Modified     time.Time      `db:"catmodified" json:"modified"`
	Items        []CategoryItem `json:"items"`
}

type CategoryItem struct {
	CategoryItemId   uint16    `db:"itemid" json:"id"`
	CategoryItemName string    `db:"itemname" json:"name"`
	CategoryItemIcon string    `db:"itemicon" json:"icon"`
	Modified         time.Time `db:"itemmodified" json:"modified"`
}
