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
	"strconv"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

var (
	categoriesMutex   sync.RWMutex
	categoriesCache   map[uint16]types.Category
	categoriesEtag    time.Time
	categoriesEtagStr string

	recipesMutex       sync.RWMutex
	recipesCache       map[uint32]*types.Recipe
	recipesEtag        time.Time
	recipesEtagStr     string
	publicRecipesCache map[uint32]*types.UserRecipeSimple
	userRecipesCache   map[uint32]userRecipesCacheItem = make(map[uint32]userRecipesCacheItem)

	recipePreparationIngredients map[int32][]*types.Ingredient

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

type userRecipesCacheItem struct {
	Etag  string
	Items *map[uint32]*types.UserRecipeSimple
}

func LoadCategories() {
	query := "SELECT * FROM categoryitemsview"
	var categories []dbCategory

	err := Db.Select(&categories, query)
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
	return categoriesCache, categoriesEtagStr
}

func LoadUnits() {
	query := "SELECT * FROM unitsview"
	var units []dbUnit

	err := Db.Select(&units, query)
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

func LoadRecipes() {
	query := "SELECT * FROM `allrecipes`"
	var recipes []DbRecipe

	err := Db.Select(&recipes, query)
	if err != nil {
		log.Fatalf("Failed to load recipes: %v", err)
	}

	// Build cache
	recipesMutex.Lock()
	recipesCache = make(map[uint32]*types.Recipe)
	publicRecipesCache = make(map[uint32]*types.UserRecipeSimple)

	for _, recipe := range recipes {

		// log.Printf("  - %d: %s / %s / %s", recipe.Id, recipe.NameDe, recipe.NameEn, recipe.NameFr)

		_, userobj := GetUser(int(recipe.UserId.Int32))

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
			Categories:     []*types.RecipeCategoryitem{},
			Preparation:    []*types.Preparation{},
			Pictures:       []*types.Picture{},
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

		recipeItem.SimpleRecipe = types.RecipeSimple{
			AiGenerated:      &recipeItem.AiGenerated,
			AiTranslatedTime: &recipeItem.AiTranslatedTime,
			Categories:       &recipeItem.Categories,
			CreatedTime:      &recipeItem.CreatedTime,
			Difficulty:       &recipeItem.Difficulty,
			EditedByUserTime: &recipeItem.EditedByUserTime,
			Id:               &recipeItem.Id,
			Localization:     &recipeItem.Localization,
			ModifiedTime:     &recipeItem.ModifiedTime,
			Pictures:         []*types.Picture{},
			PublishedTime:    &recipeItem.PublishedTime,
			ServingsCount:    &recipeItem.ServingsCount,
			SimpleStruct:     true,
			Statistics:       &recipeItem.Statistics,
			Timing:           &recipeItem.Timing,
			User:             &recipeItem.User,
			UserLocale:       &recipeItem.UserLocale,
		}

		if recipe.SharedPublic {
			publicRecipesCache[recipe.Id] = &types.UserRecipeSimple{
				RecipeSimple: &recipeItem.SimpleRecipe,
				Reason:       types.SharedPublic,
			}
		}

		if recipe.Modified.After(recipesEtag) {
			recipesEtag = recipe.Modified
		}

		recipesCache[recipe.Id] = &recipeItem

	}

	loadRecipesCategories()
	loadRecipesIngredients()
	loadRecipesPreparation()
	loadRecipesPictures()

	// log.Printf("Public recipes cache ETag: %v", recipesEtagStr)
	recipesEtagStr = hash(recipesEtag.Format(time.RFC3339) + strconv.Itoa(len(recipes)))
	recipesMutex.Unlock()

	log.Printf("Loaded %d recipes into cache", len(recipes))

	loadCollectionItems()

	go wsNotifyRecipesChanged()
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

func loadRecipesCategories() {
	query := "SELECT * FROM `recipe_categories`"
	var items []struct {
		RecipeId int       `db:"recipe_id"`
		ItemId   int       `db:"catitem_id"`
		UserId   int       `db:"user_id"`
		Created  time.Time `db:"created"`
	}

	err := Db.Select(&items, query)
	if err != nil {
		log.Fatalf("Failed to load recipes categories: %v", err)
	}

	for _, item := range items {
		// log.Printf("  - %d < %d", item.RecipeId, item.ItemId)
		_, user := GetUser(item.UserId)
		recipe, found := recipesCache[uint32(item.RecipeId)]
		if !found {
			continue
		}

		recipe.Categories = append(recipe.Categories, &types.RecipeCategoryitem{
			ItemId:  uint16(item.ItemId),
			UserId:  user.SimpleProfile,
			Created: item.Created,
		})

		if item.Created.After(recipesEtag) {
			recipesEtag = item.Created
		}
	}

	// log.Printf("Loaded %d recipes categories into cache", len(items))
}

func loadRecipesIngredients() {
	fn := "loadRecipesIngredients"

	query := "SELECT * FROM `recipe_ingredients` ORDER BY `step_id`, `sortindex`"
	var items []struct {
		Id            uint64            `db:"ingredient_id"`
		RecipeId      uint32            `db:"recipe_id"`
		StepId        types.NullInt32   `db:"step_id"`
		UnitId        types.NullInt32   `db:"unit_id"`
		Index         uint16            `db:"sortindex"`
		Quantity      types.NullFloat64 `db:"quantity"`
		DescriptionDe string            `db:"description_de"`
		DescriptionEn string            `db:"description_en"`
		DescriptionFr string            `db:"description_fr"`
	}

	err := Db.Select(&items, query)
	if err != nil {
		log.Fatalf("Failed to load recipes ingredients: %v", err)
	}

	recipePreparationIngredients = make(map[int32][]*types.Ingredient)

	for _, item := range items {
		recipe, err := GetRecipeInternal(item.RecipeId)
		if err != nil {
			log.Fatalf("%v: Invalid recipe %d for ingredient %d", fn, item.RecipeId, item.Id)
		}

		stepid := item.StepId.Int32
		if stepid == 0 {
			stepid = int32(item.RecipeId) * -1
		}
		// log.Printf("  - %d < %d: %s", stepid, item.Id, item.DescriptionDe)

		ing := &types.Ingredient{
			Id:        item.Id,
			Quantity:  item.Quantity,
			RecipeId:  recipe.Id,
			StepId:    uint64(stepid),
			SortIndex: item.Index,
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
		}

		recipePreparationIngredients[stepid] = append(recipePreparationIngredients[stepid], ing)

	}

	log.Printf("Loaded %d ingredients into cache", len(items))

}

func loadRecipesPreparation() {
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

	err := Db.Select(&steps, query)
	if err != nil {
		log.Fatalf("Failed to load recipes preparation steps: %v", err)
	}

	for _, step := range steps {
		// log.Printf("  - %d < %d: %s", step.RecipeId, step.Index, step.TitleDe)
		recipe, found := recipesCache[uint32(step.RecipeId)]
		if !found {
			continue
		}

		tempid := int32(recipe.Id) * -1

		if len(recipePreparationIngredients[tempid]) > 0 && step.Index == 0 {
			recipePreparationIngredients[int32(step.Id)] = recipePreparationIngredients[tempid]

			for _, ing := range recipePreparationIngredients[int32(step.Id)] {
				ing.StepId = step.Id
			}
		}

		step := types.Preparation{
			Id:          step.Id,
			Index:       step.Index,
			Ingredients: recipePreparationIngredients[int32(step.Id)],
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
				Preparing: ConvertTimingValue(step.Preparing),
				Cooking:   ConvertTimingValue(step.Cooking),
				Waiting:   ConvertTimingValue(step.Waiting),
			},
		}

		step.Timing.Total = ConvertTotalTimingValue(step.Timing)

		recipe.Preparation = append(recipe.Preparation, &step)
	}

	// log.Printf("Loaded %d recipes preparation steps into cache", len(steps))
}

func loadRecipesPictures() {
	query := "SELECT * FROM `recipe_pictures` WHERE `deleted` IS NULL ORDER BY `recipe_id`, `sortindex`"
	var items []dbPicture

	err := Db.Select(&items, query)
	if err != nil {
		log.Fatalf("Failed to load recipes pictures: %v", err)
	}

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

		recipe, found := recipesCache[uint32(item.RecipeId)]
		if !found {
			continue
		}

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
		AddPictureToRecipe(recipe, &picture)

	}

	// log.Printf("Loaded %d recipes categories into cache", len(items))
}

func AddPictureToRecipe(recipe *types.Recipe, picture *types.Picture) {
	recipe.Pictures = append(recipe.Pictures, picture)

	if len(recipe.Pictures) == 1 {
		recipe.SimpleRecipe.Pictures = append(recipe.SimpleRecipe.Pictures, picture)
	}

	if !picture.Dimension.Generated.Valid {
		ThumbnailGenerationRequests = append(ThumbnailGenerationRequests, PictureRequiresThumbnail{
			RecipeId:  picture.RecipeId,
			PictureId: picture.Id,
			Picture:   picture,
			Index:     picture.Index,
		})
	}

	if picture.Uploaded.After(recipesEtag) {
		recipesEtag = picture.Uploaded
	}
}

// Returns a list of recipes user is allowed to see. The returned array includes a reason for each item. This list is cached and is not updated until recipes etag changes.
func GetRecipes(user *types.UserProfile) (map[uint32]*types.UserRecipeSimple, string) {
	if user.Id == 0 {
		return publicRecipesCache, recipesEtagStr
	}

	cachemap, found := userRecipesCache[user.Id]
	if found && cachemap.Etag == recipesEtagStr {
		return *cachemap.Items, recipesEtagStr
	}

	userRecipes := make(map[uint32]*types.UserRecipeSimple)
	for _, recipe := range recipesCache {
		if recipe.SharedPublic {
			userRecipes[recipe.Id] = &types.UserRecipeSimple{
				RecipeSimple: &recipe.SimpleRecipe,
				Reason:       types.SharedPublic,
			}
			continue
		}

		if recipe.SharedInternal {
			// no user check required (user id 0 already excluded before loop)
			userRecipes[recipe.Id] = &types.UserRecipeSimple{
				RecipeSimple: &recipe.SimpleRecipe,
				Reason:       types.SharedInternal,
			}
			continue
		}

		if recipe.OwnerUserId.Valid && recipe.OwnerUserId.Int32 == int32(user.Id) {
			userRecipes[recipe.Id] = &types.UserRecipeSimple{
				RecipeSimple: &recipe.SimpleRecipe,
				Reason:       types.IsOwner,
			}
			continue
		}

		if user.Admin {
			userRecipes[recipe.Id] = &types.UserRecipeSimple{
				RecipeSimple: &recipe.SimpleRecipe,
				Reason:       types.IsAdmin,
			}
			continue
		}

	}

	userRecipesCache[user.Id] = userRecipesCacheItem{
		Etag:  recipesEtagStr,
		Items: &userRecipes,
	}

	return userRecipes, recipesEtagStr
}

func GetRecipesEtag() string {
	return recipesEtagStr
}

func getRecipeCommon(id uint32, user types.UserProfile) (*types.Recipe, error) {
	fn := fmt.Sprintf("getRecipeCommon(%d, %d, %s)", id, user.Id, user.DisplayName)

	recipe, found := recipesCache[id]
	if !found {
		log.Printf("%s: Not found in cache", fn)
		return nil, errors.New("not found")
	}

	if recipe.SharedPublic {
		// public available -> ok
		return recipe, nil
	} else if user.Id == 0 {
		// user not logged in - as not public available -> error
		log.Printf("%s: Not public shared", fn)
		return nil, errors.New("not found")
	}

	// user loggedin

	if recipe.SharedInternal {
		// internal available -> ok
		return recipe, nil
	}

	if recipe.OwnerUserId.Valid && recipe.OwnerUserId.Int32 != int32(user.Id) {
		// not shared but owner -> ok
		return recipe, nil
	}

	if user.Admin {
		// not shared but user is admin -> ok
		return recipe, nil
	}

	log.Printf("%s: No permission", fn)
	return nil, errors.New("not found")
}

func GetRecipe(id uint32, c *gin.Context) (*types.Recipe, error) {
	_, _, user, _ := GetSelf(c)
	return getRecipeCommon(id, *user)
}

func GetRecipeWs(id uint32, conn *wsConnection) (*types.Recipe, error) {
	_, _, user, _ := GetSelfByState(conn.ConnectionParams.Session)
	return getRecipeCommon(id, *user)
}

func GetRecipeInternal(id uint32) (*types.Recipe, error) {
	recipe, found := recipesCache[id]
	if !found {
		return nil, errors.New("not found")
	}
	return recipe, nil
}

func GetPicture(recipe *types.Recipe, pictureId uint32) (int, types.Picture, error) {
	for i, p := range recipe.Pictures {
		if p.Id == pictureId {
			return i, *p, nil
		}
	}
	return -1, types.Picture{}, errors.New("not found")
}

func PutRecipeLocalization(recipe *types.Recipe) (bool, error) {
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

		res, err := putRecipeLocalizationMetadata(tx, recipe, l)
		if err != nil || !res {
			_ = tx.Rollback()
			return false, err
		}

		res, err = putRecipeLocalizationPictures(tx, recipe, l)
		if err != nil || !res {
			_ = tx.Rollback()
			return false, err
		}

		res, err = putRecipeLocalizationPreparation(tx, recipe, l)
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
	touchRecipe(recipe)

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

		putRecipeLocalizationPreparationIngredients(tx, prep, lang)
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
	log.Printf("Updating recipe modified from %v to %v", recipe.ModifiedTime, time.Now())
	log.Printf("RecipesEtag = %v", recipesEtag)

	recipe.ModifiedTime = time.Now()
	if recipe.ModifiedTime.After(recipesEtag) {
		recipesEtag = recipe.ModifiedTime
	}
	recipesEtagStr = hash(recipesEtag.Format(time.RFC3339) + strconv.Itoa(len(recipesCache)))
	go wsNotifyRecipesChanged()
}

func GenerateResizedPictureVersions(recipeId uint32, pictureId uint32) (bool, error) {
	// log.Printf("Generating resized picture variants")

	recipe, err := GetRecipeInternal(recipeId)
	if err != nil {
		return false, err
	}

	i, picture, err := GetPicture(recipe, pictureId)
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

	/* recipe.Pictures[i] = &picture
	recipesMutex.Lock()
	touchRecipe(recipe)
	recipesCache[recipe.Id] = recipe
	recipesMutex.Unlock() */
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
