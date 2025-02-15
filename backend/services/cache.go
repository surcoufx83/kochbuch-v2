package services

import (
	"kochbuch-v2-backend/types"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	cacheMutex sync.RWMutex

	categoriesCache   map[string]types.Category
	categoriesEtag    time.Time
	categoriesEtagStr string

	unitsCache   map[uint8]types.Unit
	unitsEtag    time.Time
	unitsEtagStr string
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

type dbUnit struct {
	Id             uint8     `db:"unit_id"`
	ReplacedById   uint8     `db:"supersededby_unitid"`
	SavedAsId      uint8     `db:"saveas_unitid"`
	SavedAsFactor  float32   `db:"saveas_factor"`
	Localized      bool      `db:"localized"`
	SgNameDe       string    `db:"sg_name_de"`
	SgNameEn       string    `db:"sg_name_en"`
	SgNameFr       string    `db:"sg_name_fr"`
	PlNameDe       string    `db:"pl_name_de"`
	PlNameEn       string    `db:"pl_name_en"`
	PlNameFr       string    `db:"pl_name_fr"`
	DecimalPlaces  uint8     `db:"decimal_places"`
	ShowAsFraction bool      `db:"fractional"`
	Created        time.Time `db:"created"`
	Modified       time.Time `db:"updated"`
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
				Id:       category.CategoryId,
				Name:     category.CategoryName,
				Icon:     category.CategoryIcon,
				Modified: category.CategoryModified,
				Items:    []types.CategoryItem{},
			}
		}

		tempcat.Items = append(tempcat.Items, types.CategoryItem{
			Id:       category.ItemId,
			Name:     category.ItemName,
			Icon:     category.ItemIcon,
			Modified: category.ItemModified,
		})
		categoriesCache[category.CategoryName] = tempcat

		if category.CategoryModified.After(categoriesEtag) {
			categoriesEtag = category.CategoryModified
		} else if category.ItemModified.After(categoriesEtag) {
			categoriesEtag = category.ItemModified
		}
	}

	categoriesEtagStr = hash(categoriesEtag.Format(time.RFC3339) + strconv.Itoa(len(categories)))
	cacheMutex.Unlock()
	log.Printf("Loaded %d categories into cache", len(categories))
	log.Printf("Categories cache ETag: %v", categoriesEtagStr)
}

func GetCategories() (map[string]types.Category, string) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	return categoriesCache, categoriesEtagStr
}

func LoadUnits(db *sqlx.DB) {
	query := "SELECT * FROM unitsview"
	var units []dbUnit

	err := db.Select(&units, query)
	if err != nil {
		log.Fatalf("Failed to load units: %v", err)
	}

	// Build cache
	cacheMutex.Lock()
	unitsCache = make(map[uint8]types.Unit)
	for _, unit := range units {
		unitsCache[unit.Id] = types.Unit{
			Id: unit.Id,
			Localization: map[string]types.UnitLocalization{
				"de": {
					NameSingular: unit.SgNameDe,
					NamePlural:   unit.PlNameDe,
				},
				"en": {
					NameSingular: unit.SgNameEn,
					NamePlural:   unit.PlNameEn,
				},
				"fr": {
					NameSingular: unit.SgNameFr,
					NamePlural:   unit.PlNameFr,
				},
			},
			CreatedTime:    unit.Created,
			ModifiedTime:   unit.Modified,
			ReplacedById:   unit.ReplacedById,
			SavedAsId:      unit.SavedAsId,
			SavedAsFactor:  unit.SavedAsFactor,
			DecimalPlaces:  unit.DecimalPlaces,
			ShowAsFraction: unit.ShowAsFraction,
		}

		if unit.Modified.After(unitsEtag) {
			unitsEtag = unit.Modified
		}
	}
	unitsEtagStr = hash(categoriesEtag.Format(time.RFC3339) + strconv.Itoa(len(units)))
	cacheMutex.Unlock()
	log.Printf("Loaded %d units into cache", len(units))
	log.Printf("Units cache ETag: %v", unitsEtagStr)
}

func GetUnits() (map[uint8]types.Unit, string) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	return unitsCache, unitsEtagStr
}
