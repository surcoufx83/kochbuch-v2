package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kochbuch-v2-backend/types"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	AiAvailable              bool
	aiApiKey                 string
	aiTranslatorId           string
	aiTranslatorModel        string
	aiTranslatorName         string
	aiTranslatorInstructions string
	aiBaseUrl                string
)

type getAssistantsResp struct {
	ObjectType string                      `json:"object"`
	Data       []getAssistantsRespListitem `json:"data"`
}

type getAssistantsRespListitem struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	ObjectType string `json:"object"`
}

type createAssistantReq struct {
	Name         string                   `json:"name"`
	Tools        []createAssistantReqTool `json:"tools"`
	Model        string                   `json:"model"`
	Instructions string                   `json:"instructions"`
}

type createAssistantReqTool struct {
	Type string `json:"type"`
}

type recipeTranslationReq struct {
	TargetLanguage string                         `json:"targetLang"`
	Data           recipeTranslationReqRecipeData `json:"data"`
}

type recipeTranslationReqRecipeData struct {
	Recipe recipeTranslationReqRecipe `json:"recipe"`
}

type recipeTranslationReqRecipe struct {
	Title             string                            `json:"title"`
	Description       string                            `json:"description"`
	SourceDescription string                            `json:"sourceDescription"`
	Pictures          []recipeTranslationReqPicture     `json:"pictures"`
	Preparation       []recipeTranslationReqPreparation `json:"preparation"`
}

type recipeTranslationReqIngredient struct {
	Title string `json:"title"`
}

type recipeTranslationReqPicture struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type recipeTranslationReqPreparation struct {
	Title        string                           `json:"title"`
	Instructions string                           `json:"instruct"`
	Ingredients  []recipeTranslationReqIngredient `json:"ingredients"`
}

func AiConnect() (int, error) {
	aiBaseUrl = "https://api.openai.com/v1/"
	aiTranslatorName = "Kochbuch-v2 Recipe translator"
	aiTranslatorInstructions = "You are a technical assistant specialized in translating cooking recipes. The input is a JSON object with the following structure: { \"targetLang\": \"<language code>\", \"data\": { ... } }. Your task is to translate all string values within the 'data' object into the language specified by 'targetLang'. IMPORTANT: Do not change any keys or the overall structure of the JSON. The 'data' object may represent a complete recipe—including picture descriptions, preparation instructions, and ingredients—or only parts of it. Use culturally appropriate expressions for the target language (e.g., British phrases for English, German phrases for German, and French phrases for French). Do not translate picture file names if they appear as values. If any error occurs, return a JSON object with the format: { \"error\": \"message\" }."

	aiApiKey = os.Getenv("AI_APIKey")
	aiTranslatorModel = os.Getenv("AI_APIModel")
	if aiApiKey == "" {
		AiAvailable = false
		log.Println("Skipped OpenAI API integration due to missing env variable AI_APIKey")
		return http.StatusServiceUnavailable, nil
	}

	if aiTranslatorModel == "" {
		aiTranslatorModel = "gpt-4o-mini"
	}

	log.Println("Checking OpenAI API integration")
	code, err := aiGetAssistants()

	if code == http.StatusNotFound {
		code, err = aiCreateAssistant()
	}

	if code == http.StatusOK {
		AiAvailable = true
		log.Println("OpenAI API integration available")
		log.Printf("  > Translator = %v", aiTranslatorId)
		go AiAutoTranslation()
	}

	return code, err

}

func aiGetAssistants() (int, error) {
	req, err := http.NewRequest("GET", aiBaseUrl+"assistants", nil)
	if err != nil {
		log.Printf("  > Failed preparing request to list assistants: %v", err)
		return http.StatusServiceUnavailable, err
	}

	// Set required headers
	req.Header.Set("Authorization", "Bearer "+aiApiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("  > Failed sending request to list assistants: %v", err)
		return http.StatusServiceUnavailable, err
	}
	defer resp.Body.Close()

	// Check for a successful HTTP status code.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > Unexpected status to list assistants: %d  body = %s", resp.StatusCode, string(body))
		return http.StatusServiceUnavailable, err
	}

	var assistantsResp getAssistantsResp
	if err := json.NewDecoder(resp.Body).Decode(&assistantsResp); err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > Failed to parse response body = %s", string(body))
		return http.StatusServiceUnavailable, err
	}

	if len(assistantsResp.Data) == 0 {
		log.Print("  > No assistants available")
		return http.StatusNotFound, err
	}

	for _, item := range assistantsResp.Data {
		if item.Name == aiTranslatorName {
			aiTranslatorId = item.Id
			return http.StatusOK, nil
		}
	}

	return http.StatusNotFound, err
}

func aiCreateAssistant() (int, error) {
	log.Println("Creating new translator assistant")

	reqData := createAssistantReq{
		Name:         aiTranslatorName,
		Tools:        []createAssistantReqTool{},
		Model:        aiTranslatorModel,
		Instructions: aiTranslatorInstructions,
	}

	reqDataBytes, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("  > Failed preparing request payload: %v", err)
		return http.StatusServiceUnavailable, err
	}

	bodyReader := bytes.NewReader(reqDataBytes)
	req, err := http.NewRequest("POST", aiBaseUrl+"assistants", bodyReader)
	if err != nil {
		log.Printf("  > Failed preparing request to create assistant: %v", err)
		return http.StatusServiceUnavailable, err
	}

	// Set required headers
	req.Header.Set("Authorization", "Bearer "+aiApiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("  > Failed sending request to create assistant: %v", err)
		return http.StatusServiceUnavailable, err
	}
	defer resp.Body.Close()

	// Check for a successful HTTP status code.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > Unexpected status to create assistant: %d  body = %s", resp.StatusCode, string(body))
		return http.StatusServiceUnavailable, err
	}

	var assistantsResp getAssistantsRespListitem
	if err := json.NewDecoder(resp.Body).Decode(&assistantsResp); err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > Failed to parse response body = %s", string(body))
		return http.StatusServiceUnavailable, err
	}

	log.Printf("  > id = %s", assistantsResp.Id)
	aiTranslatorId = assistantsResp.Id

	return http.StatusAccepted, nil

}

func AiAutoTranslation() {
	time.Sleep(10 * time.Second)

	if !AiAvailable {
		log.Println("Skipping AI auto translation.")
		return
	}

	log.Println("Checking for entities to translate...")
	aiTranslateRecipes()
}

func aiTranslateRecipes() {
	query := "SELECT `recipe_id` FROM `recipes_translationmissing` LIMIT 1"
	var recipeids []uint32

	err := Db.Select(&recipeids, query)
	if err != nil {
		log.Printf("  > Failed to load recipes: %v", err)
		return
	}

	log.Printf("  > Query returned %v items", len(recipeids))
	for _, recipeid := range recipeids {
		recipe, err := GetRecipeInternal(recipeid)

		if err != nil {
			log.Printf("  > Failed to load recipe %d: %v", recipeid, err)
			return
		}

		t, _ := json.Marshal(recipe)
		fmt.Println(string(t))

		fmt.Println(recipe.Id == recipeid)
		fmt.Println(recipe.UserLocale == "de")

		if recipe.Id == recipeid {
			if recipe.UserLocale == "de" {
				aiTranslateRecipe(recipe, recipe.UserLocale, "en")

			}
		}

	}
}

func aiTranslateRecipe(recipe types.Recipe, from string, to string) {
	log.Printf("  > Creating translateble object from %v to %v", from, to)

	reqData := recipeTranslationReq{
		TargetLanguage: to,
		Data: recipeTranslationReqRecipeData{
			Recipe: recipeTranslationReqRecipe{
				Title:             recipe.Localization[from].Title,
				Description:       recipe.Localization[from].Description,
				SourceDescription: recipe.Localization[from].SourceDescription,
				Pictures:          []recipeTranslationReqPicture{},
				Preparation:       []recipeTranslationReqPreparation{},
			},
		},
	}

	for _, prep := range recipe.Preparation {
		step := recipeTranslationReqPreparation{
			Title:        prep.Localization[from].Title,
			Instructions: prep.Localization[from].Instructions,
			Ingredients:  []recipeTranslationReqIngredient{},
		}
		for _, ing := range prep.Ingredients {
			steping := recipeTranslationReqIngredient{
				Title: ing.Localization[from].Title,
			}
			step.Ingredients = append(step.Ingredients, steping)
		}
		reqData.Data.Recipe.Preparation = append(reqData.Data.Recipe.Preparation, step)
	}

	for _, pic := range recipe.Pictures {
		reqpic := recipeTranslationReqPicture{
			Name:        pic.Localization[from].Name,
			Description: pic.Localization[from].Description,
		}
		reqData.Data.Recipe.Pictures = append(reqData.Data.Recipe.Pictures, reqpic)
	}

	reqJson, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("  > Failed to create JSON data: %v", err)
		return
	}
	fmt.Println(string(reqJson))

}
