package types

import "time"

// Category represents a category in the database
type Category struct {
	ItemId       uint16    `db:"itemid"`
	ItemName     string    `db:"itemname"`
	ItemIcon     string    `db:"itemicon"`
	ItemModified time.Time `db:"itemmodified"`
	CatId        uint16    `db:"catid"`
	CatName      string    `db:"catname"`
	CatIcon      string    `db:"caticon"`
	CatModified  time.Time `db:"catmodified"`
}
