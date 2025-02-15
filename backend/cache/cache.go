package cache

import (
	"kochbuch-v2-backend/types"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	cacheMutex sync.RWMutex

	categoriesCache   map[uint16]types.Category
	categoriesEtag    time.Time
	categoriesEtagStr string
)

func LoadCategories(db *sqlx.DB) {
	query := "SELECT * FROM categoryitemsview"
	var categories []types.Category

	err := db.Select(&categories, query)
	if err != nil {
		log.Fatalf("Failed to load categories: %v", err)
	}

	// Build cache
	cacheMutex.Lock()
	categoriesCache = make(map[uint16]types.Category)
	for _, category := range categories {
		categoriesCache[category.ItemId] = category
		if category.ItemModified.After(categoriesEtag) {
			categoriesEtag = category.ItemModified
		} else if category.CatModified.After(categoriesEtag) {
			categoriesEtag = category.CatModified
		}
	}
	categoriesEtagStr = categoriesEtag.Format(time.RFC3339)
	cacheMutex.Unlock()
	log.Printf("Loaded %d categories into cache", len(categories))
	log.Printf("Categories cache ETag: %v", categoriesEtagStr)
}

func GetCategories() (map[uint16]types.Category, string) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	return categoriesCache, categoriesEtagStr
}
