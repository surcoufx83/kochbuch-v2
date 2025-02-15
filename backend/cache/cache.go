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

	categoriesCache   map[string]types.Category
	categoriesEtag    time.Time
	categoriesEtagStr string
)

type dbCategory struct {
	ItemId           uint16    `db:"itemid"`
	ItemName         string    `db:"itemname"`
	ItemIcon         string    `db:"itemicon"`
	ItemModified     time.Time `db:"itemmodified"`
	CategoryId       uint16    `db:"catid"`
	CategoryName     string    `db:"catname"`
	CategoryIcon     string    `db:"caticon"`
	CategoryModified time.Time `db:"catmodified"`
}

func LoadCategories(db *sqlx.DB) {
	query := "SELECT * FROM categoryitemsview"
	var categories []dbCategory

	err := db.Select(&categories, query)
	if err != nil {
		log.Fatalf("Failed to load categories: %v", err)
	}

	// Build cache
	cacheMutex.Lock()
	categoriesCache = make(map[string]types.Category)
	for _, category := range categories {

		tempcat := categoriesCache[category.CategoryName]

		if _, ok := categoriesCache[category.CategoryName]; !ok {
			tempcat = types.Category{
				CategoryId:   category.CategoryId,
				CategoryName: category.CategoryName,
				CategoryIcon: category.CategoryIcon,
				Modified:     category.CategoryModified,
				Items:        []types.CategoryItem{},
			}
		}

		tempcat.Items = append(tempcat.Items, types.CategoryItem{
			CategoryItemId:   category.ItemId,
			CategoryItemName: category.ItemName,
			CategoryItemIcon: category.ItemIcon,
			Modified:         category.ItemModified,
		})
		categoriesCache[category.CategoryName] = tempcat

		if category.CategoryModified.After(categoriesEtag) {
			categoriesEtag = category.CategoryModified
		} else if category.ItemModified.After(categoriesEtag) {
			categoriesEtag = category.ItemModified
		}
	}
	categoriesEtagStr = categoriesEtag.Format(time.RFC3339)
	cacheMutex.Unlock()
	log.Printf("Loaded %d categories into cache", len(categories))
	log.Printf("Categories cache ETag: %v", categoriesEtagStr)
}

func GetCategories() (map[string]types.Category, string) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	return categoriesCache, categoriesEtagStr
}
