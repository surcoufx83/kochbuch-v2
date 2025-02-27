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
	Error          types.NullString               `json:"error"`
}

type recipeTranslationReqRecipeData struct {
	Recipe AiRecipeTranslation `json:"recipe"`
}

type AiRecipeTranslation struct {
	Title             string                           `json:"title"`
	Description       string                           `json:"description"`
	SourceDescription string                           `json:"sourceDescription"`
	Pictures          []AiRecipeTranslationPicture     `json:"pictures"`
	Preparation       []AiRecipeTranslationPreparation `json:"preparation"`
}

type AiRecipeTranslationPicture struct {
	Id          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AiRecipeTranslationPreparation struct {
	Id           uint64                          `json:"id"`
	Title        string                          `json:"title"`
	Instructions string                          `json:"instruct"`
	Ingredients  []AiRecipeTranslationIngredient `json:"ingredients"`
}

type AiRecipeTranslationIngredient struct {
	Id    uint64 `json:"id"`
	Title string `json:"title"`
}

type createThreadReq struct {
	AssistantId string                  `json:"assistant_id"`
	Thread      createThreadReqMessages `json:"thread"`
}

type createThreadReqMessages struct {
	Messages []types.AiMessageRequest `json:"messages"`
}

type createThreadResp struct {
	RunID       string `json:"id"`
	Object      string `json:"object"`
	AssistantID string `json:"assistant_id"`
	ThreadID    string `json:"thread_id"`
	Status      string `json:"status"`
}

type getThreadMessagesResp struct {
	Object   string            `json:"object"`
	Messages []types.AiMessage `json:"data"`
}

func AiConnect() (int, error) {
	aiBaseUrl = "https://api.openai.com/v1/"
	aiTranslatorName = "Kochbuch-v2 Recipe translator"
	aiTranslatorInstructions = "You are a technical assistant specialized in translating cooking recipes. The input is a JSON object with the following structure: { \"targetLang\": \"<language code>\", \"data\": { ... } }. Your task is to translate all string values within the 'data' object into the language specified by 'targetLang'. IMPORTANT: Do not change any keys or the overall structure of the JSON. The 'data' object may represent a complete recipe—including picture descriptions, preparation instructions, and ingredients—or only parts of it. Use culturally appropriate expressions for the target language (e.g., British phrases for English, German phrases for German, and French phrases for French). Do not translate ids and picture file names if they appear as values. If any error occurs, return a JSON object with the format: { \"error\": \"message\" }."

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

	return http.StatusOK, nil

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

		if recipe.Id == recipeid {
			if recipe.UserLocale == "de" {
				aiTranslateRecipe(recipe, recipe.UserLocale, "en")
				aiTranslateRecipe(recipe, recipe.UserLocale, "fr")
			}
		}

	}
}

func aiTranslateRecipe(recipe types.Recipe, from string, to string) {
	log.Printf("  > %v > Creating translateble object from %v to %v", to, from, to)

	reqData := recipeTranslationReq{
		TargetLanguage: to,
		Data: recipeTranslationReqRecipeData{
			Recipe: AiRecipeTranslation{
				Title:             recipe.Localization[from].Title,
				Description:       recipe.Localization[from].Description,
				SourceDescription: recipe.Localization[from].SourceDescription,
				Pictures:          []AiRecipeTranslationPicture{},
				Preparation:       []AiRecipeTranslationPreparation{},
			},
		},
	}

	for _, prep := range recipe.Preparation {
		step := AiRecipeTranslationPreparation{
			Id:           prep.Id,
			Title:        prep.Localization[from].Title,
			Instructions: prep.Localization[from].Instructions,
			Ingredients:  []AiRecipeTranslationIngredient{},
		}
		for _, ing := range prep.Ingredients {
			steping := AiRecipeTranslationIngredient{
				Id:    ing.Id,
				Title: ing.Localization[from].Title,
			}
			step.Ingredients = append(step.Ingredients, steping)
		}
		reqData.Data.Recipe.Preparation = append(reqData.Data.Recipe.Preparation, step)
	}

	for _, pic := range recipe.Pictures {
		reqpic := AiRecipeTranslationPicture{
			Id:          pic.Id,
			Name:        pic.Localization[from].Name,
			Description: pic.Localization[from].Description,
		}
		reqData.Data.Recipe.Pictures = append(reqData.Data.Recipe.Pictures, reqpic)
	}

	reqJson, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("  > %v > Failed to create JSON data: %v", to, err)
		return
	}

	res, resp := aiCreateThread(reqJson, to)
	if !res {
		return
	}

	go aiFetchThreadResponse(resp, recipe, to)

}

func aiCreateThread(message []byte, to string) (bool, createThreadResp) {
	log.Printf("  > %v > Creating messaging thread at OpenAI API", to)

	msg := types.AiMessageRequest{
		Role:    "user",
		Content: string(message),
	}

	threadReq := createThreadReq{
		AssistantId: aiTranslatorId,
		Thread: createThreadReqMessages{
			Messages: []types.AiMessageRequest{},
		},
	}
	threadReq.Thread.Messages = append(threadReq.Thread.Messages, msg)

	reqDataBytes, err := json.Marshal(threadReq)
	if err != nil {
		log.Printf("  > %v > Failed preparing request payload: %v", to, err)
		return false, createThreadResp{}
	}

	bodyReader := bytes.NewReader(reqDataBytes)
	req, err := http.NewRequest("POST", aiBaseUrl+"threads/runs", bodyReader)
	if err != nil {
		log.Printf("  > %v > Failed preparing request to create thread: %v", to, err)
		return false, createThreadResp{}
	}

	// Set required headers
	req.Header.Set("Authorization", "Bearer "+aiApiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("  > %v > Failed sending request to create thread: %v", to, err)
		return false, createThreadResp{}
	}
	defer resp.Body.Close()

	// Check for a successful HTTP status code.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > %v > Unexpected status to create thread: %d  body = %s", to, resp.StatusCode, string(body))
		return false, createThreadResp{}
	}

	var threadResp createThreadResp
	if err := json.NewDecoder(resp.Body).Decode(&threadResp); err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > %v > Failed to parse response body = %s", to, string(body))
		return false, createThreadResp{}
	}

	return true, threadResp
}

func aiFetchThreadResponse(thread createThreadResp, recipe types.Recipe, to string) {
	log.Printf("  > %v > Waiting for translated recipe a few seconds", to)

	time.Sleep(10 * time.Second)
	log.Printf("  > %v > Querying translated recipe", to)
	log.Printf("    > %v > thread_id %v", to, thread.ThreadID)
	log.Printf("    > %v >    run_id %v", to, thread.RunID)

	req, err := http.NewRequest("GET", aiBaseUrl+"threads/"+thread.ThreadID+"/messages", nil)
	if err != nil {
		log.Printf("  > %v > Failed preparing request to get messages: %v", to, err)
		return
	}

	// Set required headers
	req.Header.Set("Authorization", "Bearer "+aiApiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("  > %v > Failed sending request to get messages: %v", to, err)
		return
	}
	defer resp.Body.Close()

	// Check for a successful HTTP status code.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > %v > Unexpected status to get messages: %d  body = %s", to, resp.StatusCode, string(body))
		return
	}

	var msgResp getThreadMessagesResp
	if err := json.NewDecoder(resp.Body).Decode(&msgResp); err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > %v > Failed to parse response body = %s", to, string(body))
		return
	}

	msgBytes, err := json.Marshal(msgResp)
	if err != nil {
		log.Printf("  > %v > Failed to JSON decode response: %v", to, err)
		return
	}

	if len(msgResp.Messages) == 0 {
		log.Printf("  > %v > Empty messages array: %v", to, err)
		fmt.Println(string(msgBytes))
		return
	}

	if msgResp.Messages[0].Role != "assistant" {
		log.Printf("  > %v > Message 0 not of expected role assistant: %v", to, err)
		fmt.Println(string(msgBytes))
		return
	}

	if len(msgResp.Messages[0].Content) == 0 {
		log.Printf("  > %v > Empty messages content array: %v", to, err)
		fmt.Println(string(msgBytes))
		return
	}

	var reply recipeTranslationReq
	if err = json.Unmarshal([]byte(msgResp.Messages[0].Content[0].Text.Value), &reply); err != nil {
		log.Printf("  > %v > Failed to JSON decode first message content: %v", to, err)
		fmt.Println(msgResp.Messages[0].Content[0].Text.Value)
		return
	}

	if reply.Error.Valid && reply.Error.String != "" {
		log.Printf("  > %v > Reply contains error message: %v", to, reply.Error.String)
		return
	}

	if reply.TargetLanguage != to {
		log.Printf("  > %v > Reply language does not match expected one: %v !- %v", to, reply.TargetLanguage, to)
		return
	}

	originRecipe, err := GetRecipeInternal(recipe.Id)
	if err != nil {
		log.Printf("  > %v > Failed getting Recipe from cache: %v", to, err)
		return
	}

	PutRecipeLocalization(originRecipe, to, reply.Data.Recipe, true)

}
