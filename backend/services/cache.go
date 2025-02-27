package services

import (
	"database/sql"
	"errors"
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

	recipePreparationIngredients map[uint64][]types.Ingredient

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
	WaitingTime            types.NullInt32  `db:"waiting_time"`
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
	query := "SELECT * FROM `allrecipes`"
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

		_, userobj := GetUser(int(recipe.PictureUserId.Int32))

		recipeItem := types.Recipe{
			Id:               recipe.Id,
			SimpleStruct:     true,
			IsFork:           false,
			OriginalRecipeId: 0,
			IsPlaceholder:    recipe.IsPlaceholder,
			SourceUrl:        recipe.SourceUrl,
			OwnerUserId:      recipe.UserId,
			LastEditUserId:   recipe.EditUserId,
			User:             userobj.SimpleProfile,
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
				Preparing: recipe.PreparingTime,
				Cooking:   recipe.CookingTime,
				Waiting:   recipe.WaitingTime,
			},
		}

		if recipe.PictureId.Valid {

			_, userobj := GetUser(int(recipe.PictureUserId.Int32))

			picture := types.Picture{
				Id:       uint32(recipe.PictureId.Int32),
				RecipeId: recipe.Id,
				User:     userobj.SimpleProfile,
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

	loadRecipesCategories(db)
	loadRecipesIngredients(db)
	loadRecipesPreparation(db)
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
		log.Printf("  - %d < %d", item.RecipeId, item.ItemId)
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
	log.Printf("Loaded %d recipes categories into cache", len(items))
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
		log.Printf("  - %d < %d: %s", item.StepId, item.Id, item.DescriptionDe)
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
	log.Printf("Loaded %d recipes preparation steps into cache", len(items))
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
		log.Printf("  - %d < %d: %s", step.RecipeId, step.Index, step.TitleDe)
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
	log.Printf("Loaded %d recipes preparation steps into cache", len(steps))
}

func GetRecipes(c *gin.Context) (map[uint32]types.Recipe, string) {
	code, _, user, err := GetSelf(c)

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

func PutRecipeLocalization(recipe types.Recipe, lang string, locale AiRecipeTranslation, ai bool) (bool, error) {
	log.Printf("Updating recipe localisation %v %v", recipe.Id, recipe.Localization[lang].Title)

	tx, err := Db.Begin()
	if err != nil {
		log.Printf("  > Failed starting transaction: %v", err)
		return false, err
	}

	res, err := putRecipeLocalizationMetadata(tx, &recipe, lang, &locale, ai)
	if err != nil || !res {
		_ = tx.Rollback()
		return false, err
	}

	res, err = putRecipeLocalizationPictures(tx, &recipe, lang, &locale)
	if err != nil || !res {
		_ = tx.Rollback()
		return false, err
	}

	res, err = putRecipeLocalizationPreparation(tx, &recipe, lang, &locale)
	if err != nil || !res {
		_ = tx.Rollback()
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	go LoadRecipes(Db)

	return true, nil
}

func putRecipeLocalizationMetadata(tx *sql.Tx, recipe *types.Recipe, lang string, locale *AiRecipeTranslation, ai bool) (bool, error) {
	log.Printf("  > Patching general data")
	stmt, err := tx.Prepare("UPDATE `recipes` SET `localized` = ?, `name_" + lang + "` = ?, `description_" + lang + "` = ?, `source_description_" + lang + "` = ? WHERE `recipe_id` = ?")
	if err != nil {
		log.Printf("  > Failed preparing stmt: %v", err)
		return false, err
	}

	if !ai {
		ai = recipe.AiLocalized
	}

	_, err = stmt.Exec(ai, locale.Title, locale.Description, locale.SourceDescription, recipe.Id)
	if err != nil {
		log.Printf("  > Failed executing stmt: %v", err)
		return false, err
	}

	return true, nil
}

func putRecipeLocalizationPictures(tx *sql.Tx, recipe *types.Recipe, lang string, locale *AiRecipeTranslation) (bool, error) {
	log.Printf("  > Patching pictures")
	stmt, err := tx.Prepare("UPDATE `recipe_pictures` SET `name_" + lang + "` = ?, `description_" + lang + "` = ? WHERE `picture_id` = ?")
	if err != nil {
		log.Printf("  > Failed preparing stmt: %v", err)
		return false, err
	}

	for _, pic := range recipe.Pictures {
		for _, langpic := range locale.Pictures {
			if pic.Id != langpic.Id {
				continue
			}

			_, err = stmt.Exec(langpic.Name, langpic.Description, pic.Id)
			if err != nil {
				log.Printf("  > Failed executing stmt: %v", err)
				return false, err
			}

		}
	}

	return true, nil
}

func putRecipeLocalizationPreparation(tx *sql.Tx, recipe *types.Recipe, lang string, locale *AiRecipeTranslation) (bool, error) {
	log.Printf("  > Patching preparation steps")
	stmt, err := tx.Prepare("UPDATE `recipe_steps` SET `title_" + lang + "` = ?, `instruct_" + lang + "` = ? WHERE `step_id` = ?")
	if err != nil {
		log.Printf("  > Failed preparing stmt: %v", err)
		return false, err
	}

	for _, prep := range recipe.Preparation {
		for _, langprep := range locale.Preparation {
			if prep.Id != langprep.Id {
				continue
			}

			_, err = stmt.Exec(langprep.Title, langprep.Instructions, prep.Id)
			if err != nil {
				log.Printf("  > Failed executing stmt: %v", err)
				return false, err
			}

			putRecipeLocalizationPreparationIngredients(tx, &prep, lang, &langprep)

		}
	}

	return true, nil
}

func putRecipeLocalizationPreparationIngredients(tx *sql.Tx, prep *types.Preparation, lang string, locale *AiRecipeTranslationPreparation) (bool, error) {
	log.Printf("    > Patching ingredients for a step")
	stmt, err := tx.Prepare("UPDATE `recipe_ingredients` SET `description_" + lang + "` = ? WHERE `ingredient_id` = ?")
	if err != nil {
		log.Printf("  > Failed preparing stmt: %v", err)
		return false, err
	}

	for _, ing := range prep.Ingredients {
		for _, langing := range locale.Ingredients {
			if ing.Id != langing.Id {
				continue
			}

			_, err = stmt.Exec(langing.Title, ing.Id)
			if err != nil {
				log.Printf("  > Failed executing stmt: %v", err)
				return false, err
			}

		}
	}

	return true, nil
}
