package types

import (
	"slices"
	"time"
)

var (
	QueryCollectionAddItem string = "INSERT INTO `user_collection_recipes`(`collection_id`, `recipe_id`, `is_owner`, `remarks`) VALUES(?, ?, ?, ?)"

	QueryCollectionRemoveItem string = "DELETE FROM `user_collection_recipes` WHERE `collection_id` = ? AND `recipe_id` = ? LIMIT 1"
)

type Collection struct {
	Created     time.Time `json:"created" db:"created"`
	Deleted     NullTime  `json:"deleted" db:"deleted"`
	Description string    `json:"description" db:"description"`
	Id          uint32    `json:"id" db:"collection_id"`
	Modified    time.Time `json:"modified" db:"modified"`
	Name        string    `json:"title" db:"title"`
	Published   NullTime  `json:"published" db:"published"`
	UserId      uint32    `json:"-" db:"user_id"`

	Items []*CollectionItem `json:"items" db:"-"`
}

type CollectionItem struct {
	CollectionId uint32    `json:"-" db:"collection_id"`
	Created      time.Time `json:"created" db:"created"`
	IsOwner      bool      `json:"isOwner" db:"is_owner"`
	Modified     time.Time `json:"modified" db:"modified"`
	RecipeId     uint32    `json:"recipeId" db:"recipe_id"`
	Remarks      string    `json:"remarks" db:"remarks"`
	UserId       uint32    `json:"-" db:"user_id"`
}

func (c *Collection) AddItem(u *UserProfile, r *Recipe, remarks string) {
	c.Items = append(c.Items, &CollectionItem{
		CollectionId: c.Id,
		Created:      time.Now(),
		IsOwner:      u.Id == uint32(r.OwnerUserId.Int32),
		Modified:     time.Now(),
		RecipeId:     r.Id,
		Remarks:      remarks,
		UserId:       u.Id,
	})
	c.Modified = time.Now()
}

func (c *Collection) Contains(r *Recipe) bool {
	for _, i := range c.Items {
		if i.RecipeId == r.Id {
			return true
		}
	}
	return false
}

func (c *Collection) RemoveItem(r *Recipe) {
	removeat := -1
	for i, item := range c.Items {
		if item.RecipeId == r.Id {
			removeat = i
			break
		}
	}
	if removeat == -1 {
		return
	}
	c.Items = slices.Replace(c.Items, removeat, removeat+1)
}
