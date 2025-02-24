package services

import (
	"kochbuch-v2-backend/types"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	categoriesMutex   sync.RWMutex
	categoriesCache   map[uint16]types.Category
	categoriesEtag    time.Time
	categoriesEtagStr string

	recipesMutex       sync.RWMutex
	recipesCache       map[uint32]types.Recipe
	recipesEtag        time.Time
	recipesEtagStr     string
	publicRecipesCache map[uint32]types.Recipe

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
	Id                     uint32           `db:"recipe_id"`
	UserId                 types.NullInt32  `db:"user_id"`
	EditUserId             types.NullInt32  `db:"edit_user_id"`
	AiGenerated            bool             `db:"aigenerated"`
	AiLocalized            bool             `db:"localized"`
	IsPlaceholder          bool             `db:"placeholder"`
	SharedInternal         bool             `db:"shared_internal"`
	SharedPublic           bool             `db:"shared_external"`
	Locale                 string           `db:"locale"`
	NameDe                 string           `db:"name_de"`
	NameEn                 string           `db:"name_en"`
	NameFr                 string           `db:"name_fr"`
	DescriptionDe          string           `db:"description_de"`
	DescriptionEn          string           `db:"description_en"`
	DescriptionFr          string           `db:"description_fr"`
	ServingsCount          uint8            `db:"servings_count"`
	SourceDescriptionDe    string           `db:"source_description_de"`
	SourceDescriptionEn    string           `db:"source_description_en"`
	SourceDescriptionFr    string           `db:"source_description_fr"`
	SourceUrl              string           `db:"source_url"`
	Created                time.Time        `db:"created"`
	Modified               time.Time        `db:"modified"`
	Published              types.NullTime   `db:"published"`
	Difficulty             uint8            `db:"difficulty"`
	IngredientsGroupByStep bool             `db:"ingredientsGroupByStep"`
	PictureId              types.NullInt32  `db:"picture_id"`
	PictureUserId          types.NullInt32  `db:"picture_user_id"`
	PictureIndex           types.NullInt32  `db:"picture_sortindex"`
	PictureNameDe          types.NullString `db:"picture_name_de"`
	PictureNameEn          types.NullString `db:"picture_name_en"`
	PictureNameFr          types.NullString `db:"picture_name_fr"`
	PictureDescriptionDe   types.NullString `db:"picture_description_de"`
	PictureDescriptionEn   types.NullString `db:"picture_description_en"`
	PictureDescriptionFr   types.NullString `db:"picture_description_fr"`
	PictureFilename        types.NullString `db:"picture_filename"`
	PictureFullPath        types.NullString `db:"picture_fullpath"`
	PictureUploaded        types.NullTime   `db:"picture_uploaded"`
	PictureWidth           types.NullInt32  `db:"picture_width"`
	PictureHeight          types.NullInt32  `db:"picture_height"`
	ViewsCount             uint32           `db:"views"`
	CookedCount            uint32           `db:"cooked"`
	VotesCount             uint32           `db:"votes"`
	VotesSum               uint32           `db:"votesum"`
	VotesAvg               float32          `db:"avgvotes"`
	RatingsCount           uint32           `db:"ratings"`
	RatingsSum             uint32           `db:"ratesum"`
	RatingsAvg             float32          `db:"avgratings"`
	StepsCount             uint8            `db:"stepscount"`
	PreparingTime          types.NullInt32  `db:"preparing_time"`
	CookingTime            types.NullInt32  `db:"cooking_time"`
	WitingTime             types.NullInt32  `db:"waiting_time"`
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

func LoadRecipes(db *sqlx.DB) {
	query := "SELECT * FROM allrecipes"
	var recipes []DbRecipe

	err := db.Select(&recipes, query)
	if err != nil {
		log.Fatalf("Failed to load recipes: %v", err)
	}

	// Build cache
	recipesMutex.Lock()
	recipesCache = make(map[uint32]types.Recipe)
	publicRecipesCache = make(map[uint32]types.Recipe)
	for _, recipe := range recipes {

		log.Printf("  - %d: %s / %s / %s", recipe.Id, recipe.NameDe, recipe.NameEn, recipe.NameFr)

		recipeItem := types.Recipe{
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
			Preparation:    []types.Preparation{},
			Pictures:       []types.Picture{},
			ServingsCount:  recipe.ServingsCount,
			Difficulty:     recipe.Difficulty,
			SharedInternal: recipe.SharedInternal,
			SharedPublic:   recipe.SharedPublic,
			CreatedTime:    recipe.Created,
			ModifiedTime:   recipe.Modified,
			PublishedTime:  recipe.Published,
		}

		if recipe.PictureId.Valid {

			picture := types.Picture{
				Id:       uint32(recipe.PictureId.Int32),
				RecipeId: recipe.Id,
				UserId:   recipe.PictureUserId,
				Index:    uint8(recipe.PictureIndex.Int32),
				Localization: map[string]types.PictureLocalization{
					"de": {
						Name:        recipe.PictureNameDe.String,
						Description: recipe.PictureDescriptionDe.String,
					},
					"en": {
						Name:        recipe.PictureNameEn.String,
						Description: recipe.PictureDescriptionEn.String,
					},
					"fr": {
						Name:        recipe.PictureNameFr.String,
						Description: recipe.PictureDescriptionFr.String,
					},
				},
				FileName: recipe.PictureFilename.String,
				FullPath: recipe.PictureFullPath.String,
				Uploaded: recipe.PictureUploaded,
				Dimension: types.PictureDimension{
					Height: recipe.PictureHeight.Int32,
					Width:  recipe.PictureWidth.Int32,
				},
			}
			recipeItem.Pictures = append(recipeItem.Pictures, picture)

		}

		if recipe.Modified.After(recipesEtag) {
			recipesEtag = recipe.Modified
		}

		recipesCache[recipe.Id] = recipeItem

		if recipe.SharedPublic {
			publicRecipesCache[recipe.Id] = recipesCache[recipe.Id]
		}

	}

	recipesEtagStr = hash(recipesEtag.Format(time.RFC3339) + strconv.Itoa(len(recipes)))
	recipesMutex.Unlock()
	log.Printf("Loaded %d recipes into cache", len(recipes))
	log.Printf("Public recipes cache ETag: %v", recipesEtagStr)
}

func GetRecipes(c *gin.Context) (map[uint32]types.Recipe, string) {
	code, user, err := GetSelf(c)

	if err != nil || code != http.StatusOK || user.Id == 0 {
		return publicRecipesCache, recipesEtagStr
	}

	recipesMutex.RLock()
	defer recipesMutex.RUnlock()

	userRecipes := make(map[uint32]types.Recipe)
	for _, recipe := range recipesCache {
		if recipe.SharedPublic || recipe.SharedInternal || (recipe.OwnerUserId.Valid && recipe.OwnerUserId.Int32 == int32(user.Id)) {
			userRecipes[recipe.Id] = recipe
		}
	}
	return userRecipes, recipesEtagStr
}

func GetRecipesEtag() string {
	return recipesEtagStr
}
