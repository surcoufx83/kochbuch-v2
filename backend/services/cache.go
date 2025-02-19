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
	categoriesMutex   sync.RWMutex
	categoriesCache   map[uint16]types.Category
	categoriesEtag    time.Time
	categoriesEtagStr string

	publicRecipesMutex   sync.RWMutex
	publicRecipesCache   map[uint32]types.Recipe
	publicRecipesEtag    time.Time
	publicRecipesEtagStr string

	unitsMutex   sync.RWMutex
	unitsCache   map[uint8]types.Unit
	unitsEtag    time.Time
	unitsEtagStr string
)

type dbCategory struct {
	ItemId           uint16    `db:"item_id"`
	ItemNameDe       string    `db:"item_name_de"`
	ItemNameEn       string    `db:"item_name_en"`
	ItemNameFr       string    `db:"item_name_fr"`
	ItemIcon         string    `db:"item_icon"`
	ItemModified     time.Time `db:"item_modified"`
	CategoryId       uint16    `db:"cat_id"`
	CategoryNameDe   string    `db:"cat_name_de"`
	CategoryNameEn   string    `db:"cat_name_en"`
	CategoryNameFr   string    `db:"cat_name_fr"`
	CategoryIcon     string    `db:"cat_icon"`
	CategoryModified time.Time `db:"cat_modified"`
}

type DbRecipe struct {
	Id                     uint32    `db:"recipe_id"`
	UserId                 uint32    `db:"user_id"`
	EditUserId             uint32    `db:"edit_user_id"`
	AiGenerated            bool      `db:"aigenerated"`
	AiLocalized            bool      `db:"localized"`
	IsPlaceholder          bool      `db:"placeholder"`
	SharedInternal         bool      `db:"shared_internal"`
	SharedPublic           bool      `db:"shared_external"`
	Locale                 string    `db:"locale"`
	NameDe                 string    `db:"name_de"`
	NameEn                 string    `db:"name_en"`
	NameFr                 string    `db:"name_fr"`
	DescriptionDe          string    `db:"description_de"`
	DescriptionEn          string    `db:"description_en"`
	DescriptionFr          string    `db:"description_fr"`
	ServingsCount          uint8     `db:"servings_count"`
	SourceDescriptionDe    string    `db:"source_description_de"`
	SourceDescriptionEn    string    `db:"source_description_en"`
	SourceDescriptionFr    string    `db:"source_description_fr"`
	SourceUrl              string    `db:"source_url"`
	Created                time.Time `db:"created"`
	Modified               time.Time `db:"modified"`
	Published              time.Time `db:"published"`
	Difficulty             uint8     `db:"difficulty"`
	IngredientsGroupByStep bool      `db:"ingredientsGroupByStep"`
	PictureId              uint32    `db:"picture_id"`
	PictureIndex           uint8     `db:"picture_sortindex"`
	PictureName            string    `db:"picture_name"`
	PictureDescription     string    `db:"picture_description"`
	PictureHash            string    `db:"picture_hash"`
	PictureFilename        string    `db:"picture_filename"`
	PictureFullPath        string    `db:"picture_full_path"`
	PictureUploaded        time.Time `db:"picture_uploaded"`
	PictureWidth           uint16    `db:"picture_width"`
	PictureHeight          uint16    `db:"picture_height"`
	ViewsCount             uint32    `db:"views"`
	CookedCount            uint32    `db:"cooked"`
	VotesCount             uint32    `db:"votes"`
	VotesSum               uint32    `db:"votesum"`
	VotesAvg               float32   `db:"avgvotes"`
	RatingsCount           uint32    `db:"ratings"`
	RatingsSum             uint32    `db:"ratesum"`
	RatingsAvg             float32   `db:"avgratings"`
	StepsCount             uint8     `db:"stepscount"`
	PreparationTime        int16     `db:"preparationtime"`
	CookingTime            int16     `db:"cookingtime"`
	ChillTime              int16     `db:"chilltime"`
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
	categoriesMutex.Lock()
	categoriesCache = make(map[uint16]types.Category)
	for _, category := range categories {

		tempcat := categoriesCache[category.CategoryId]

		if _, ok := categoriesCache[category.CategoryId]; !ok {
			tempcat = types.Category{
				Id: category.CategoryId,
				Localization: map[string]types.NameLocalization{
					"de": {
						Name: category.CategoryNameDe,
					},
					"en": {
						Name: category.CategoryNameEn,
					},
					"fr": {
						Name: category.CategoryNameFr,
					},
				},
				Icon:     category.CategoryIcon,
				Modified: category.CategoryModified,
				Items:    make(map[uint16]types.CategoryItem),
			}
		}

		tempcat.Items[category.ItemId] = types.CategoryItem{
			Id: category.ItemId,
			Localization: map[string]types.NameLocalization{
				"de": {
					Name: category.ItemNameDe,
				},
				"en": {
					Name: category.ItemNameEn,
				},
				"fr": {
					Name: category.ItemNameFr,
				},
			},
			Icon:     category.ItemIcon,
			Modified: category.ItemModified,
		}
		categoriesCache[category.CategoryId] = tempcat

		if category.CategoryModified.After(categoriesEtag) {
			categoriesEtag = category.CategoryModified
		} else if category.ItemModified.After(categoriesEtag) {
			categoriesEtag = category.ItemModified
		}
	}

	categoriesEtagStr = hash(categoriesEtag.Format(time.RFC3339) + strconv.Itoa(len(categories)))
	categoriesMutex.Unlock()
	log.Printf("Loaded %d categories into cache", len(categories))
	log.Printf("Categories cache ETag: %v", categoriesEtagStr)
}

func GetCategories() (map[uint16]types.Category, string) {
	categoriesMutex.RLock()
	defer categoriesMutex.RUnlock()

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
	unitsMutex.Lock()
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
	unitsMutex.Unlock()
	log.Printf("Loaded %d units into cache", len(units))
	log.Printf("Units cache ETag: %v", unitsEtagStr)
}

func GetUnits() (map[uint8]types.Unit, string) {
	unitsMutex.RLock()
	defer unitsMutex.RUnlock()

	return unitsCache, unitsEtagStr
}

func LoadPublicRecipes(db *sqlx.DB) {
	query := "SELECT * FROM allrecipes_nouser"
	var recipes []DbRecipe

	err := db.Select(&recipes, query)
	if err != nil {
		log.Fatalf("Failed to load recipes: %v", err)
	}

	// Build cache
	publicRecipesMutex.Lock()
	publicRecipesCache = make(map[uint32]types.Recipe)
	for _, recipe := range recipes {

		publicRecipesCache[recipe.Id] = types.Recipe{
			Id:               recipe.Id,
			SimpleStruct:     true,
			IsFork:           false,
			OriginalRecipeId: 0,
			IsPlaceholder:    recipe.IsPlaceholder,
			SourceUrl:        recipe.SourceUrl,
			OwnerUserId:      recipe.UserId,
			LastEditUserId:   recipe.EditUserId,
			AiGenerated:      recipe.AiGenerated,
			AiLocalized:      recipe.AiLocalized,
			UserLocale:       recipe.Locale,
			Localization: map[string]types.RecipeLocalization{
				"de": {
					Title:             recipe.NameDe,
					Description:       recipe.DescriptionDe,
					SourceDescription: recipe.SourceDescriptionDe,
				},
				"en": {
					Title:             recipe.NameEn,
					Description:       recipe.DescriptionEn,
					SourceDescription: recipe.SourceDescriptionEn,
				},
				"fr": {
					Title:             recipe.NameFr,
					Description:       recipe.DescriptionFr,
					SourceDescription: recipe.SourceDescriptionFr,
				},
			},
			Categories:     []uint16{},
			ServingsCount:  recipe.ServingsCount,
			Difficulty:     recipe.Difficulty,
			SharedInternal: recipe.SharedInternal,
			SharedPublic:   recipe.SharedPublic,
			CreatedTime:    recipe.Created,
			ModifiedTime:   recipe.Modified,
			PublishedTime:  recipe.Published,
		}

		if recipe.Modified.After(publicRecipesEtag) {
			publicRecipesEtag = recipe.Modified
		}
	}

	publicRecipesEtagStr = hash(publicRecipesEtag.Format(time.RFC3339) + strconv.Itoa(len(recipes)))
	publicRecipesMutex.Unlock()
	log.Printf("Loaded %d recipes into cache", len(recipes))
	log.Printf("Public recipes cache ETag: %v", publicRecipesEtagStr)
}

func GetPublicRecipes() (map[uint32]types.Recipe, string) {
	publicRecipesMutex.RLock()
	defer publicRecipesMutex.RUnlock()

	return publicRecipesCache, publicRecipesEtagStr
}
