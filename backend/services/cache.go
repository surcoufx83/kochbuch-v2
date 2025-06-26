package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"kochbuch-v2-backend/types"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
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
	publicRecipesCache map[uint32]*types.RecipeSimple

	recipePreparationIngredients map[uint64][]types.Ingredient

	unitsMutex   sync.RWMutex
	unitsCache   map[uint8]types.Unit
	unitsEtag    time.Time
	unitsEtagStr string

	Locales []string

	ThumbnailSizes []int
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
	Id                     uint32          `db:"recipe_id"`
	UserId                 types.NullInt32 `db:"user_id"`
	EditUserId             types.NullInt32 `db:"edit_user_id"`
	AiGenerated            bool            `db:"aigenerated"`
	AiTranslatedTime       types.NullTime  `db:"localized"`
	EditedByUserTime       types.NullTime  `db:"edited"`
	IsPlaceholder          bool            `db:"placeholder"`
	SharedInternal         bool            `db:"shared_internal"`
	SharedPublic           bool            `db:"shared_external"`
	Locale                 string          `db:"locale"`
	NameDe                 string          `db:"name_de"`
	NameEn                 string          `db:"name_en"`
	NameFr                 string          `db:"name_fr"`
	DescriptionDe          string          `db:"description_de"`
	DescriptionEn          string          `db:"description_en"`
	DescriptionFr          string          `db:"description_fr"`
	ServingsCount          uint8           `db:"servings_count"`
	SourceDescriptionDe    string          `db:"source_description_de"`
	SourceDescriptionEn    string          `db:"source_description_en"`
	SourceDescriptionFr    string          `db:"source_description_fr"`
	SourceUrl              string          `db:"source_url"`
	Created                time.Time       `db:"created"`
	Modified               time.Time       `db:"modified"`
	Published              types.NullTime  `db:"published"`
	Difficulty             uint8           `db:"difficulty"`
	IngredientsGroupByStep bool            `db:"ingredientsGroupByStep"`
	ViewsCount             uint32          `db:"views"`
	CookedCount            uint32          `db:"cooked"`
	VotesCount             uint32          `db:"votes"`
	VotesSum               uint32          `db:"votesum"`
	VotesAvg               float32         `db:"avgvotes"`
	RatingsCount           uint32          `db:"ratings"`
	RatingsSum             uint32          `db:"ratesum"`
	RatingsAvg             float32         `db:"avgratings"`
	StepsCount             uint8           `db:"stepscount"`
	PreparingTime          types.NullInt32 `db:"preparing_time"`
	CookingTime            types.NullInt32 `db:"cooking_time"`
	WaitingTime            types.NullInt32 `db:"waiting_time"`
}

type dbPicture struct {
	PictureId     uint32           `db:"picture_id"`
	RecipeId      uint32           `db:"recipe_id"`
	UserId        sql.NullInt32    `db:"user_id"`
	SortIndex     uint8            `db:"sortindex"`
	NameDe        string           `db:"name_de"`
	NameEn        string           `db:"name_en"`
	NameFr        string           `db:"name_fr"`
	DescriptionDe string           `db:"description_de"`
	DescriptionEn string           `db:"description_en"`
	DescriptionFr string           `db:"description_fr"`
	Hash          string           `db:"hash"`
	Filename      string           `db:"filename"`
	Fullpath      string           `db:"fullpath"`
	Uploaded      time.Time        `db:"uploaded"`
	Deleted       types.NullTime   `db:"deleted"`
	Width         uint16           `db:"width"`
	Height        uint16           `db:"height"`
	ThbSizes      types.NullString `db:"thb_sizes"`
	ThbGenerated  types.NullTime   `db:"thb_generated"`
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
	// log.Printf("Categories cache ETag: %v", categoriesEtagStr)
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
	// log.Printf("Units cache ETag: %v", unitsEtagStr)
}

func GetUnits() (map[uint8]types.Unit, string) {
	unitsMutex.RLock()
	defer unitsMutex.RUnlock()

	return unitsCache, unitsEtagStr
}

func LoadRecipes(db *sqlx.DB) {
	query := "SELECT * FROM `allrecipes`"
	var recipes []DbRecipe

	err := db.Select(&recipes, query)
	if err != nil {
		log.Fatalf("Failed to load recipes: %v", err)
	}

	// Build cache
	recipesMutex.Lock()
	recipesCache = make(map[uint32]types.Recipe)
	publicRecipesCache = make(map[uint32]*types.RecipeSimple)
	for _, recipe := range recipes {

		// log.Printf("  - %d: %s / %s / %s", recipe.Id, recipe.NameDe, recipe.NameEn, recipe.NameFr)

		_, userobj := GetUser(int(recipe.EditUserId.Int32))

		recipeItem := types.Recipe{
			Id:               recipe.Id,
			SimpleStruct:     false,
			IsFork:           false,
			OriginalRecipeId: 0,
			IsPlaceholder:    recipe.IsPlaceholder,
			SourceUrl:        recipe.SourceUrl,
			OwnerUserId:      recipe.UserId,
			LastEditUserId:   recipe.EditUserId,
			User:             userobj.SimpleProfile,
			AiGenerated:      recipe.AiGenerated,
			AiTranslatedTime: recipe.AiTranslatedTime,
			EditedByUserTime: recipe.EditedByUserTime,
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
			Categories:     []types.RecipeCategoryitem{},
			Preparation:    []types.Preparation{},
			Pictures:       []types.Picture{},
			ServingsCount:  recipe.ServingsCount,
			Difficulty:     recipe.Difficulty,
			SharedInternal: recipe.SharedInternal,
			SharedPublic:   recipe.SharedPublic,
			CreatedTime:    recipe.Created,
			ModifiedTime:   recipe.Modified,
			PublishedTime:  recipe.Published,
			Statistics: types.RecipeStatistics{
				Steps:  recipe.StepsCount,
				Views:  recipe.ViewsCount,
				Cooked: recipe.CookedCount,
				Votes: types.RecipeStatisticsItem{
					Avg:   recipe.VotesAvg,
					Count: recipe.VotesCount,
				},
				Ratings: types.RecipeStatisticsItem{
					Avg:   recipe.RatingsAvg,
					Count: recipe.RatingsCount,
				},
			},
			Timing: types.PreparationTiming{
				Preparing: ConvertTimingValue(recipe.PreparingTime),
				Cooking:   ConvertTimingValue(recipe.CookingTime),
				Waiting:   ConvertTimingValue(recipe.WaitingTime),
			},
		}

		recipeItem.Timing.Total = ConvertTotalTimingValue(recipeItem.Timing)

		if recipe.Modified.After(recipesEtag) {
			recipesEtag = recipe.Modified
		}

		recipesCache[recipe.Id] = recipeItem

	}

	recipesEtagStr = hash(recipesEtag.Format(time.RFC3339) + strconv.Itoa(len(recipes)))
	recipesMutex.Unlock()
	// log.Printf("Public recipes cache ETag: %v", recipesEtagStr)

	loadRecipesCategories(db)
	loadRecipesIngredients(db)
	loadRecipesPreparation(db)
	loadRecipesPictures(db)

	for _, recipe := range recipesCache {
		if recipe.SharedPublic {
			publicRecipesCache[recipe.Id] = ConvertToRecipeSimple(&recipe)
		}
	}

	log.Printf("Loaded %d recipes into cache", len(recipes))
	wsNotifyRecipesChanged()
}

func ConvertTimingValue(dbvalue types.NullInt32) types.NullInt32 {
	if !dbvalue.Valid {
		return dbvalue
	}
	if dbvalue.Int32 < 0 {
		return types.NullInt32{
			Valid: false,
		}
	}
	return dbvalue
}

func ConvertTotalTimingValue(dbvalue types.PreparationTiming) types.NullInt32 {
	if !dbvalue.Cooking.Valid && !dbvalue.Preparing.Valid && !dbvalue.Waiting.Valid {
		return dbvalue.Cooking
	}
	result := types.NullInt32{
		Valid: true,
		Int32: 0,
	}
	if dbvalue.Cooking.Valid && dbvalue.Cooking.Int32 > 0 {
		result.Int32 += dbvalue.Cooking.Int32
	}
	if dbvalue.Preparing.Valid && dbvalue.Preparing.Int32 > 0 {
		result.Int32 += dbvalue.Preparing.Int32
	}
	if dbvalue.Waiting.Valid && dbvalue.Waiting.Int32 > 0 {
		result.Int32 += dbvalue.Waiting.Int32
	}
	return result
}

func ConvertToRecipeSimple(r *types.Recipe) *types.RecipeSimple {
	var firstPicture []types.Picture
	if len(r.Pictures) > 0 {
		firstPicture = []types.Picture{r.Pictures[0]}
	}

	return &types.RecipeSimple{
		Id:               r.Id,
		SimpleStruct:     true,
		User:             r.User,
		UserLocale:       r.UserLocale,
		Localization:     r.Localization,
		Categories:       r.Categories,
		Pictures:         firstPicture,
		ServingsCount:    r.ServingsCount,
		Difficulty:       r.Difficulty,
		Statistics:       r.Statistics,
		Timing:           r.Timing,
		AiTranslatedTime: r.AiTranslatedTime,
		CreatedTime:      r.CreatedTime,
		EditedByUserTime: r.EditedByUserTime,
		ModifiedTime:     r.ModifiedTime,
		PublishedTime:    r.PublishedTime,
	}
}

func loadRecipesCategories(db *sqlx.DB) {
	query := "SELECT * FROM `recipe_categories`"
	var items []struct {
		RecipeId int       `db:"recipe_id"`
		ItemId   int       `db:"catitem_id"`
		UserId   int       `db:"user_id"`
		Created  time.Time `db:"created"`
	}

	err := db.Select(&items, query)
	if err != nil {
		log.Fatalf("Failed to load recipes categories: %v", err)
	}

	recipesMutex.Lock()
	for _, item := range items {
		// log.Printf("  - %d < %d", item.RecipeId, item.ItemId)
		_, user := GetUser(item.UserId)
		recipe := recipesCache[uint32(item.RecipeId)]
		recipe.Categories = append(recipe.Categories, types.RecipeCategoryitem{
			ItemId:  uint16(item.ItemId),
			UserId:  user.SimpleProfile,
			Created: item.Created,
		})
		recipesCache[uint32(item.RecipeId)] = recipe
	}

	recipesMutex.Unlock()
	// log.Printf("Loaded %d recipes categories into cache", len(items))
}

func loadRecipesIngredients(db *sqlx.DB) {
	query := "SELECT * FROM `recipe_ingredients` WHERE `step_id` IS NOT NULL ORDER BY `step_id`, `sortindex`"
	var items []struct {
		Id            uint64            `db:"ingredient_id"`
		RecipeId      uint32            `db:"recipe_id"`
		StepId        uint64            `db:"step_id"`
		UnitId        types.NullInt32   `db:"unit_id"`
		Index         uint16            `db:"sortindex"`
		Quantity      types.NullFloat64 `db:"quantity"`
		DescriptionDe string            `db:"description_de"`
		DescriptionEn string            `db:"description_en"`
		DescriptionFr string            `db:"description_fr"`
	}

	err := db.Select(&items, query)
	if err != nil {
		log.Fatalf("Failed to load recipes ingredients: %v", err)
	}

	recipePreparationIngredients = make(map[uint64][]types.Ingredient)

	recipesMutex.Lock()
	for _, item := range items {
		// log.Printf("  - %d < %d: %s", item.StepId, item.Id, item.DescriptionDe)
		recipePreparationIngredients[item.StepId] = append(recipePreparationIngredients[item.StepId], types.Ingredient{
			Id:        item.Id,
			SortIndex: item.Index,
			Quantity:  item.Quantity,
			UnitId:    item.UnitId,
			Localization: map[string]types.IngredientLocalization{
				"de": {
					Title: item.DescriptionDe,
				},
				"en": {
					Title: item.DescriptionEn,
				},
				"fr": {
					Title: item.DescriptionFr,
				},
			},
		})
	}

	recipesMutex.Unlock()
	// log.Printf("Loaded %d recipes preparation steps into cache", len(items))
}

func loadRecipesPreparation(db *sqlx.DB) {
	query := "SELECT * FROM `recipe_steps` ORDER BY `recipe_id`, `sortindex`"
	var steps []struct {
		Id         uint64          `db:"step_id"`
		RecipeId   uint32          `db:"recipe_id"`
		Index      uint8           `db:"sortindex"`
		TitleDe    string          `db:"title_de"`
		TitleEn    string          `db:"title_en"`
		TitleFr    string          `db:"title_fr"`
		InstructDe string          `db:"instruct_de"`
		InstructEn string          `db:"instruct_en"`
		InstructFr string          `db:"instruct_fr"`
		Preparing  types.NullInt32 `db:"preparing"`
		Cooking    types.NullInt32 `db:"cooking"`
		Waiting    types.NullInt32 `db:"waiting"`
	}

	err := db.Select(&steps, query)
	if err != nil {
		log.Fatalf("Failed to load recipes preparation steps: %v", err)
	}

	recipesMutex.Lock()
	for _, step := range steps {
		// log.Printf("  - %d < %d: %s", step.RecipeId, step.Index, step.TitleDe)
		recipe := recipesCache[uint32(step.RecipeId)]
		recipe.Preparation = append(recipe.Preparation, types.Preparation{
			Id:          step.Id,
			Index:       step.Index,
			Ingredients: recipePreparationIngredients[step.Id],
			Localization: map[string]types.PreparationLocalization{
				"de": {
					Title:        step.TitleDe,
					Instructions: step.InstructDe,
				},
				"en": {
					Title:        step.TitleEn,
					Instructions: step.InstructEn,
				},
				"fr": {
					Title:        step.TitleFr,
					Instructions: step.InstructFr,
				},
			},
			Timing: types.PreparationTiming{
				Preparing: step.Preparing,
				Cooking:   step.Cooking,
				Waiting:   step.Waiting,
			},
		})
		recipesCache[uint32(step.RecipeId)] = recipe
	}

	recipesMutex.Unlock()
	// log.Printf("Loaded %d recipes preparation steps into cache", len(steps))
}

func loadRecipesPictures(db *sqlx.DB) {
	query := "SELECT * FROM `recipe_pictures` WHERE `deleted` IS NULL ORDER BY `recipe_id`, `sortindex`"
	var items []dbPicture

	err := db.Select(&items, query)
	if err != nil {
		log.Fatalf("Failed to load recipes pictures: %v", err)
	}

	recipesMutex.Lock()
	for _, item := range items {
		// log.Printf("  - %d < %d", item.RecipeId, item.ItemId)
		_, user := GetUser(int(item.UserId.Int32))
		thbsizes := []int{}

		if item.ThbSizes.Valid {
			if err = json.Unmarshal([]byte(item.ThbSizes.String), &thbsizes); err != nil {
				item.ThbSizes = types.NullString{
					Valid:  true,
					String: "[]",
				}
				item.ThbGenerated = types.NullTime{
					Valid: false,
				}
			}
		}

		recipe := recipesCache[uint32(item.RecipeId)]
		basename, ext := GetBasenameAndExtension(item.Filename)

		picture := types.Picture{
			Id:       item.PictureId,
			RecipeId: item.RecipeId,
			User:     user.SimpleProfile,
			Index:    uint8(len(recipe.Pictures)),
			Localization: map[string]types.PictureLocalization{
				"de": {
					Name:        item.NameDe,
					Description: item.DescriptionDe,
				},
				"en": {
					Name:        item.NameEn,
					Description: item.DescriptionEn,
				},
				"fr": {
					Name:        item.NameFr,
					Description: item.DescriptionFr,
				},
			},
			FileName: item.Filename,
			BaseName: basename,
			Ext:      ext,
			FullPath: item.Fullpath,
			Uploaded: item.Uploaded,
			Dimension: types.PictureDimension{
				Height:         int(item.Height),
				Width:          int(item.Width),
				GeneratedSizes: thbsizes,
				Generated:      item.ThbGenerated,
			},
		}
		recipe.Pictures = append(recipe.Pictures, picture)
		recipesCache[uint32(item.RecipeId)] = recipe

		if item.ThbGenerated.Valid {
			for _, size := range ThumbnailSizes {
				if !slices.Contains(thbsizes, size) {
					item.ThbSizes = types.NullString{
						Valid:  true,
						String: "[]",
					}
					item.ThbGenerated = types.NullTime{
						Valid: false,
					}
				}
			}
		}

		if !item.ThbGenerated.Valid {
			ThumbnailGenerationRequests = append(ThumbnailGenerationRequests, PictureRequiresThumbnail{
				RecipeId:  picture.RecipeId,
				PictureId: picture.Id,
				Picture:   &picture,
				Index:     picture.Index,
			})
		}

	}

	recipesMutex.Unlock()
	// log.Printf("Loaded %d recipes categories into cache", len(items))
}

func GetRecipes(user *types.UserProfile) (map[uint32]*types.RecipeSimple, string) {
	if user.Id == 0 {
		return publicRecipesCache, recipesEtagStr
	}

	recipesMutex.RLock()
	defer recipesMutex.RUnlock()

	userRecipes := make(map[uint32]*types.RecipeSimple)
	for _, recipe := range recipesCache {
		if recipe.SharedPublic || recipe.SharedInternal || (recipe.OwnerUserId.Valid && recipe.OwnerUserId.Int32 == int32(user.Id)) {
			userRecipes[recipe.Id] = ConvertToRecipeSimple(&recipe)
		}
	}
	return userRecipes, recipesEtagStr
}

func GetRecipesEtag() string {
	return recipesEtagStr
}

func GetRecipe(id uint32, c *gin.Context) (types.Recipe, error) {
	_, _, user, _ := GetSelf(c)

	recipesMutex.RLock()
	defer recipesMutex.RUnlock()

	for _, recipe := range recipesCache {
		if recipe.Id != id {
			continue
		}

		if recipe.SharedPublic || recipe.SharedInternal || (recipe.OwnerUserId.Valid && recipe.OwnerUserId.Int32 == int32(user.Id)) {
			return recipe, nil
		}
	}
	return types.Recipe{}, errors.New("not found")
}

func GetRecipeWs(id uint32, conn *wsConnection) (types.Recipe, error) {
	_, _, user, _ := GetSelfByState(conn.ConnectionParams.Session)

	recipesMutex.RLock()
	defer recipesMutex.RUnlock()

	for _, recipe := range recipesCache {
		if recipe.Id != id {
			continue
		}

		if recipe.SharedPublic || recipe.SharedInternal || (recipe.OwnerUserId.Valid && recipe.OwnerUserId.Int32 == int32(user.Id)) {
			return recipe, nil
		}
	}
	return types.Recipe{}, errors.New("not found")
}

func GetRecipeInternal(id uint32) (types.Recipe, error) {
	recipesMutex.RLock()
	defer recipesMutex.RUnlock()

	for _, recipe := range recipesCache {
		if recipe.Id != id {
			continue
		}
		return recipe, nil
	}
	return types.Recipe{}, errors.New("not found")
}

func GetPicture(recipe *types.Recipe, pictureId uint32) (int, types.Picture, error) {
	for i, p := range recipe.Pictures {
		if p.Id == pictureId {
			return i, p, nil
		}
	}
	return -1, types.Picture{}, errors.New("not found")
}

func PutRecipeLocalization(recipe types.Recipe) (bool, error) {
	log.Printf("Updating recipe localisation %v %v", recipe.Id, recipe.Localization[recipe.UserLocale].Title)

	tx, err := Db.Begin()
	if err != nil {
		log.Printf("  > Failed starting transaction: %v", err)
		return false, err
	}

	for _, l := range Locales {

		if l == recipe.UserLocale {
			continue
		}

		res, err := putRecipeLocalizationMetadata(tx, &recipe, l)
		if err != nil || !res {
			_ = tx.Rollback()
			return false, err
		}

		res, err = putRecipeLocalizationPictures(tx, &recipe, l)
		if err != nil || !res {
			_ = tx.Rollback()
			return false, err
		}

		res, err = putRecipeLocalizationPreparation(tx, &recipe, l)
		if err != nil || !res {
			_ = tx.Rollback()
			return false, err
		}
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	log.Printf("  > Translation finished")
	go LoadRecipes(Db)

	return true, nil
}

func putRecipeLocalizationMetadata(tx *sql.Tx, recipe *types.Recipe, lang string) (bool, error) {
	log.Printf("  > Patching general data")
	stmt, err := tx.Prepare("UPDATE `recipes` SET `localized` = current_timestamp(), `name_" + lang + "` = ?, `description_" + lang + "` = ?, `source_description_" + lang + "` = ? WHERE `recipe_id` = ?")
	if err != nil {
		log.Printf("  > Failed preparing stmt: %v", err)
		return false, err
	}

	_, err = stmt.Exec(recipe.Localization[lang].Title, recipe.Localization[lang].Description, recipe.Localization[lang].SourceDescription, recipe.Id)
	if err != nil {
		log.Printf("  > Failed executing stmt: %v", err)
		return false, err
	}

	return true, nil
}

func putRecipeLocalizationPictures(tx *sql.Tx, recipe *types.Recipe, lang string) (bool, error) {
	log.Printf("  > Patching pictures")
	stmt, err := tx.Prepare("UPDATE `recipe_pictures` SET `name_" + lang + "` = ?, `description_" + lang + "` = ? WHERE `picture_id` = ?")
	if err != nil {
		log.Printf("  > Failed preparing stmt: %v", err)
		return false, err
	}

	for _, pic := range recipe.Pictures {
		_, err = stmt.Exec(pic.Localization[lang].Name, pic.Localization[lang].Description, pic.Id)
		if err != nil {
			log.Printf("  > Failed executing stmt: %v", err)
			return false, err
		}
	}

	return true, nil
}

func putRecipeLocalizationPreparation(tx *sql.Tx, recipe *types.Recipe, lang string) (bool, error) {
	log.Printf("  > Patching preparation steps")
	stmt, err := tx.Prepare("UPDATE `recipe_steps` SET `title_" + lang + "` = ?, `instruct_" + lang + "` = ? WHERE `step_id` = ?")
	if err != nil {
		log.Printf("  > Failed preparing stmt: %v", err)
		return false, err
	}

	for _, prep := range recipe.Preparation {
		_, err = stmt.Exec(prep.Localization[lang].Title, prep.Localization[lang].Instructions, prep.Id)
		if err != nil {
			log.Printf("  > Failed executing stmt: %v", err)
			return false, err
		}

		putRecipeLocalizationPreparationIngredients(tx, &prep, lang)
	}

	return true, nil
}

func putRecipeLocalizationPreparationIngredients(tx *sql.Tx, prep *types.Preparation, lang string) (bool, error) {
	log.Printf("    > Patching ingredients for a step")
	stmt, err := tx.Prepare("UPDATE `recipe_ingredients` SET `description_" + lang + "` = ? WHERE `ingredient_id` = ?")
	if err != nil {
		log.Printf("  > Failed preparing stmt: %v", err)
		return false, err
	}

	for _, ing := range prep.Ingredients {
		_, err = stmt.Exec(ing.Localization[lang].Title, ing.Id)
		if err != nil {
			log.Printf("  > Failed executing stmt: %v", err)
			return false, err
		}
	}

	return true, nil
}

func touchRecipe(recipe *types.Recipe) {
	// log.Println("Updating recipe timestamp")
	recipe.ModifiedTime = time.Now()
	if recipe.ModifiedTime.After(recipesEtag) {
		recipesEtag = recipe.ModifiedTime
		recipesEtagStr = hash(recipesEtag.Format(time.RFC3339) + strconv.Itoa(len(recipesCache)))
	}
}

func GenerateResizedPictureVersions(recipeId uint32, pictureId uint32) (bool, error) {
	// log.Printf("Generating resized picture variants")

	recipe, err := GetRecipeInternal(recipeId)
	if err != nil {
		return false, err
	}

	i, picture, err := GetPicture(&recipe, pictureId)
	if err != nil || i == -1 || picture.Id == 0 {
		return false, err
	}

	if v, err := getPictureExistsOnDisk(&picture); !v || err != nil {
		return false, err
	}

	for _, size := range ThumbnailSizes {
		res, err := generateResizedPictureVersion(&picture, size)
		if !res || err != nil {
			return false, err
		}
	}
	// log.Printf("  > %v thumbnails created", picture.FullPath)

	sizesJson, err := json.Marshal(picture.Dimension.GeneratedSizes)
	if err != nil {
		log.Printf("  > %v failed marshaling sizes: %v", picture.FullPath, err)
		return false, err
	}

	query := "UPDATE `recipe_pictures` SET `width` = ?, `height` = ?, `thb_sizes` = ?, `thb_generated` = current_timestamp() WHERE `picture_id` = ?"
	stmt, err := Db.Prepare(query)
	if err != nil {
		log.Printf("  > %v failed preparing stmt: %v", picture.FullPath, err)
		return false, err
	}

	_, err = stmt.Exec(picture.Dimension.Width, picture.Dimension.Height, string(sizesJson), picture.Id)
	if err != nil {
		log.Printf("  > Failed executing stmt: %v", err)
		return false, err
	}

	recipe.Pictures[i] = picture
	recipesMutex.Lock()
	touchRecipe(&recipe)
	recipesCache[recipe.Id] = recipe
	recipesMutex.Unlock()
	log.Printf("  > %v saved to database and updated in cache", picture.FullPath)

	return false, err

}

func getPictureExistsOnDisk(picture *types.Picture) (bool, error) {
	_, err := os.Stat(picture.FullPath)
	if err != nil {
		log.Printf("  > %v failed: %v", picture.FullPath, err)
		return false, err
	}
	return true, nil
}

func generateResizedPictureVersion(picture *types.Picture, size int) (bool, error) {
	// log.Printf("  > %v -> %d", picture.FullPath, size)

	folder := filepath.Dir(picture.FullPath)
	basename, ext := GetBasenameAndExtension(picture.FileName)

	// log.Printf("  >  > Folder = %v", folder)
	// log.Printf("  >  > Base name = %v", basename)
	// log.Printf("  >  > Extension = %v", ext)

	// open image
	imgFile, err := os.Open(picture.FullPath)
	if err != nil {
		log.Printf("  > %v failed: %v", picture.FullPath, err)
		return false, err
	}
	defer imgFile.Close()

	// get picture exif data
	exifData, err := exif.Decode(imgFile)
	if err != nil {
		log.Printf("  > %v failed reading Exif data: %v", picture.FullPath, err)
	}
	imgFile.Seek(0, 0)

	// decode image to check dimensions
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Printf("  > %v failed reading image: %v", picture.FullPath, err)
		return false, err
	}

	picture.Dimension = types.PictureDimension{
		Height:         img.Bounds().Dy(),
		Width:          img.Bounds().Dx(),
		GeneratedSizes: picture.Dimension.GeneratedSizes,
		Generated:      picture.Dimension.Generated,
	}
	// log.Printf("  >  > Size = %dx%d", picture.Dimension.Width, picture.Dimension.Height)

	// Apply EXIF orientation if available
	if exifData != nil {
		orientation, err := exifData.Get(exif.Orientation)
		if err == nil {
			orientationVal, err := orientation.Int(0)
			if err == nil {
				// log.Printf("  >  > Exif Orientation: %d", orientationVal)
				img = applyExifOrientation(img, orientationVal)
			}
		}
	}

	// created resized copy
	resizedImg := resize.Resize(uint(size), 0, img, resize.Lanczos3)

	resizedFilename := filepath.Join(folder, fmt.Sprintf("%s_%d%s", basename, size, ext))
	// log.Printf("  >  > Save as = %v", resizedFilename)
	resizedFile, err := os.Create(resizedFilename)
	if err != nil {
		log.Printf("  > %v failed: %v", picture.FullPath, err)
		return false, err
	}
	defer resizedFile.Close()

	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(resizedFile, resizedImg, nil)
	case ".png":
		err = png.Encode(resizedFile, resizedImg)
	}
	if err != nil {
		log.Printf("  > %v failed: %v", picture.FullPath, err)
		return false, err
	}

	picture.Dimension = types.PictureDimension{
		Height:         picture.Dimension.Height,
		Width:          picture.Dimension.Width,
		GeneratedSizes: append(picture.Dimension.GeneratedSizes, size),
		Generated:      picture.Dimension.Generated,
	}

	return true, nil

}

// applyExifOrientation applies the EXIF orientation to the image.
func applyExifOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 3:
		return rotate180(img)
	case 6:
		return rotate90(img)
	case 8:
		return rotate270(img)
	default:
		return img
	}
}

// rotate90 rotates the image 90 degrees clockwise.
func rotate90(img image.Image) image.Image {
	// Rotate 90 degrees clockwise
	return imaging.Rotate270(img)
}

// rotate180 rotates the image 180 degrees.
func rotate180(img image.Image) image.Image {
	// Rotate 180 degrees
	return imaging.Rotate180(img)
}

// rotate270 rotates the image 270 degrees clockwise.
func rotate270(img image.Image) image.Image {
	// Rotate 270 degrees clockwise
	return imaging.Rotate90(img)
}
