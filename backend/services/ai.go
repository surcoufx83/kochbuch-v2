package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kochbuch-v2-backend/types"
	"log"
	"net/http"
	"os"
	"slices"
	"time"
)

var (
	AiAvailable              bool
	aiApiKey                 string
	aiTranslatorId           string
	aiTranslatorModel        string
	aiTranslatorName         string
	aiTranslatorInstructions string
	aiTranslatorWaitInterval time.Duration
	aiTranslatorWaitRetries  int
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
	Name           string                   `json:"name"`
	Tools          []createAssistantReqTool `json:"tools"`
	Model          string                   `json:"model"`
	Instructions   string                   `json:"instructions"`
	ResponseFormat createAssistantReqTool   `json:"response_format"`
}

type createAssistantReqTool struct {
	Type string `json:"type"`
}

type recipeTranslationReq struct {
	From    string   `json:"from"`
	To      string   `json:"to"`
	Phrases []string `json:"phrases"`
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

type getRunStatusResp struct {
	Status string `json:"status"`
}

func AiConnect() (int, error) {
	aiBaseUrl = "https://api.openai.com/v1/"
	aiTranslatorName = "Kochbuch-v2 Recipe translator"
	aiTranslatorInstructions = "You will receive requests in the form of a JSON object containing the properties 'from' and 'to' to define the origin and target languages, along with an array of strings under the key 'phrases' that should be translated. Translate each phrase using culturally appropriate expressions for the target language (e.g., British English expressions for English, standard German for German, and common French idioms for French). Do not translate any content that represents IDs, file names, or clearly non-translatable keys/values (for example, strings that are identifiers or file references). If in doubt, leave the content unchanged. Output the result as a JSON array containing the translated phrases in the same order as provided, and nothing else. If the input format is invalid, return an appropriate error message in JSON format."

	aiTranslatorWaitInterval = 3 * time.Second
	aiTranslatorWaitRetries = 20

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
		go aiAutoTranslation()
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
		ResponseFormat: createAssistantReqTool{
			Type: "json_object",
		},
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

	// log.Printf("  > id = %s", assistantsResp.Id)
	aiTranslatorId = assistantsResp.Id

	return http.StatusOK, nil

}

func aiAutoTranslation() {
	if !AiAvailable {
		log.Println("Skipping AI auto translation.")
		return
	}

	aiTranslateRecipes()
}

func aiTranslateRecipes() {
	for {
		if aiTranslateRecipesLoop() {
			time.Sleep(5 * time.Second)
		} else {
			// log.Println("No recipes for translation")
			time.Sleep(5 * time.Minute)
		}
	}
}

func aiTranslateRecipesLoop() bool {
	query := "SELECT `recipe_id` FROM `recipes_translationmissing` LIMIT 1"
	var recipeids []uint32

	err := Db.Select(&recipeids, query)
	if err != nil {
		log.Printf("  > Failed to load recipes: %v", err)
		return false
	}

	if len(recipeids) == 0 {
		return false
	}

	// log.Printf("  > Query returned %v items", len(recipeids))
	for _, recipeid := range recipeids {
		recipe, err := GetRecipeInternal(recipeid)

		if err != nil {
			log.Printf("  > Failed to load recipe %d: %v", recipeid, err)
			return false
		}

		translations := make(map[string][]string)
		var source []string
		var tl []string

		if recipe.Id == recipeid {
			log.Printf("Translating recipe %v", recipe.Localization[recipe.UserLocale].Title)

			for _, l := range Locales {
				if recipe.UserLocale != l {
					source, tl, err = aiTranslateRecipe(&recipe, recipe.UserLocale, l)
					if err != nil {
						log.Println("  > Cancelled translation for this recipe")
						return false
					}

					translations[l] = tl
				}
			}

		}

		aiApplyTranslations(recipe, source, translations)
	}
	return true
}

func aiTranslateRecipe(recipe *types.Recipe, from string, to string) ([]string, []string, error) {
	log.Printf("  > Creating translation table %v to %v", from, to)

	// Create a string array containing all strings to translate
	// without any blank strings or duplicates
	translationTable := aiRecipeFillTranslationTable([]string{}, recipe, from)

	// Create message content
	reqData := recipeTranslationReq{
		From:    from,
		To:      to,
		Phrases: translationTable,
	}

	// Convert payload to JSON
	reqJson, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("  > %v > Failed to create JSON data: %v", to, err)
		return []string{}, []string{}, err
	}
	// fmt.Println(string(reqJson))

	// Send request to Open AI API to create and run a messaging thread
	res, createdThread := aiCreateThread(reqJson)
	if !res {
		return []string{}, []string{}, err
	}

	// Wait for the completion which may take some seconds
	res, err = aiWaitThread(createdThread)
	if !res || err != nil {
		return []string{}, []string{}, err
	}

	msg, err := aiGetThreadMessage(createdThread)
	if err != nil {
		return []string{}, []string{}, err
	}

	translatedTable, err := aiGetTranslationFromMessage(msg)
	return translationTable, translatedTable, err

}

func aiRecipeFillTranslationTable(table []string, recipe *types.Recipe, from string) []string {
	table = aiRecipeFillTranslationTableItem(table, recipe.Localization[from].Title)
	table = aiRecipeFillTranslationTableItem(table, recipe.Localization[from].Description)
	table = aiRecipeFillTranslationTableItem(table, recipe.Localization[from].SourceDescription)

	for _, p := range recipe.Pictures {
		table = aiRecipeFillTranslationTableItem(table, p.Localization[from].Name)
		table = aiRecipeFillTranslationTableItem(table, p.Localization[from].Description)
	}

	for _, p := range recipe.Preparation {
		table = aiRecipeFillTranslationTableItem(table, p.Localization[from].Title)
		table = aiRecipeFillTranslationTableItem(table, p.Localization[from].Instructions)
		for _, i := range p.Ingredients {
			table = aiRecipeFillTranslationTableItem(table, i.Localization[from].Title)
		}
	}

	return table
}

func aiRecipeFillTranslationTableItem(table []string, phrase string) []string {
	if phrase != "" && !slices.Contains(table, phrase) {
		table = append(table, phrase)
	}
	return table
}

func aiCreateThread(message []byte) (bool, createThreadResp) {
	log.Printf("  > Creating messaging thread")

	msg := types.AiMessageRequest{
		Role:    "user",
		Content: "```json\n" + string(message) + "\n```",
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
		log.Printf("  > > Failed preparing request payload: %v", err)
		return false, createThreadResp{}
	}

	bodyReader := bytes.NewReader(reqDataBytes)
	req, err := http.NewRequest("POST", aiBaseUrl+"threads/runs", bodyReader)
	if err != nil {
		log.Printf("   > > Failed preparing request to create thread: %v", err)
		return false, createThreadResp{}
	}

	req.Header.Set("Authorization", "Bearer "+aiApiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("  > > Failed sending request to create thread: %v", err)
		return false, createThreadResp{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > > Unexpected status to create thread: %d  body = %s", resp.StatusCode, string(body))
		return false, createThreadResp{}
	}

	var threadResp createThreadResp
	if err := json.NewDecoder(resp.Body).Decode(&threadResp); err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > > Failed to parse response body = %s", string(body))
		return false, createThreadResp{}
	}

	log.Printf("  > > Messaging thread created: %v", threadResp.RunID)

	return true, threadResp
}

func aiWaitThread(createdThread createThreadResp) (bool, error) {
	log.Println("  > Waiting for message thread finished")

	for i := 1; i <= aiTranslatorWaitRetries; i++ {
		time.Sleep(aiTranslatorWaitInterval)
		status, err := aiGetThreadStatus(createdThread)
		if err != nil {
			return false, err
		}
		log.Printf("  > > [%d/%d] %v", i, aiTranslatorWaitRetries, status.Status)
		if status.Status == "completed" {
			break
		}
	}

	return true, nil

}

func aiGetThreadStatus(createdThread createThreadResp) (getRunStatusResp, error) {
	log.Println("  > Requesting thread status")

	req, err := http.NewRequest("GET", aiBaseUrl+"threads/"+createdThread.ThreadID+"/runs/"+createdThread.RunID, nil)
	if err != nil {
		log.Printf("  > > Failed preparing request: %v", err)
		return getRunStatusResp{}, err
	}

	req.Header.Set("Authorization", "Bearer "+aiApiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("  > > Failed sending request: %v", err)
		return getRunStatusResp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > > Unexpected status: %d  body = %s", resp.StatusCode, string(body))
		return getRunStatusResp{}, err
	}

	var payload getRunStatusResp
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > > Failed to parse response body = %s", string(body))
		return getRunStatusResp{}, err
	}

	return payload, nil
}

func aiGetThreadMessage(createdThread createThreadResp) (types.AiMessage, error) {
	log.Println("  > Requesting assistant message")

	req, err := http.NewRequest("GET", aiBaseUrl+"threads/"+createdThread.ThreadID+"/messages", nil)
	if err != nil {
		log.Printf("  > > Failed preparing request: %v", err)
		return types.AiMessage{}, err
	}

	req.Header.Set("Authorization", "Bearer "+aiApiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("  > > Failed sending request: %v", err)
		return types.AiMessage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > > Unexpected status: %d  body = %s", resp.StatusCode, string(body))
		return types.AiMessage{}, err
	}

	var msgResp getThreadMessagesResp
	if err := json.NewDecoder(resp.Body).Decode(&msgResp); err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > > Failed to parse response body = %s", string(body))
		return types.AiMessage{}, err
	}

	if len(msgResp.Messages) == 0 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("  > > Response does not contain any messages = %s", string(body))
		return types.AiMessage{}, err
	}

	for i, msg := range msgResp.Messages {
		if i == 0 && msg.Role != "assistant" {
			fmt.Printf("  > > First message not expected role assistant: %v", msg)
		}
		if msg.Role == "assistant" {
			return msg, nil
		}
	}
	fmt.Printf("  > > No message containing role assistant! %v", msgResp)

	return types.AiMessage{}, errors.New("not found")
}

func aiGetTranslationFromMessage(msg types.AiMessage) ([]string, error) {
	log.Println("  > Getting translation from message")

	if len(msg.Content) == 0 {
		fmt.Printf("  > > Message content is empty: %v", msg)
		return []string{}, errors.New("content is empty")
	}

	if msg.Content[0].Text.Value == "" {
		fmt.Printf("  > > Message content value is empty: %v", msg)
		return []string{}, errors.New("content value is empty")
	}

	var translatedStrings []string
	if err := json.Unmarshal([]byte(msg.Content[0].Text.Value), &translatedStrings); err != nil {
		fmt.Printf("  > > Failed to decode message content text value: %v", msg.Content[0].Text.Value)
		return []string{}, errors.New("failed to decode")
	}

	return translatedStrings, nil

}

func aiApplyTranslations(recipe types.Recipe, sourcePhrases []string, translations map[string][]string) {
	log.Println("  > Update recipe with new translations")
	for _, l := range Locales {
		if recipe.UserLocale == l || len(translations[l]) == 0 {
			continue
		}
		if len(sourcePhrases) != len(translations[l]) {
			log.Printf("  > Missmatch between requested phrases and translated phrases! %d != %d", len(sourcePhrases), len(translations[l]))
			return
		}
		recipe = aiApplyTranslation(recipe, l, sourcePhrases, translations[l])
	}
	PutRecipeLocalization(recipe)
}

func aiApplyTranslation(recipe types.Recipe, l string, sourcePhrases []string, translations []string) types.Recipe {
	log.Printf("  > Update recipe with new %v translations", l)
	org := recipe.UserLocale

	recipe.Localization[l] = types.RecipeLocalization{
		Title:             aiApplyPhrase(recipe.Localization[org].Title, sourcePhrases, translations),
		Description:       aiApplyPhrase(recipe.Localization[org].Description, sourcePhrases, translations),
		SourceDescription: aiApplyPhrase(recipe.Localization[org].SourceDescription, sourcePhrases, translations),
	}

	for i, p := range recipe.Pictures {
		recipe.Pictures[i].Localization[l] = types.PictureLocalization{
			Name:        aiApplyPhrase(p.Localization[l].Name, sourcePhrases, translations),
			Description: aiApplyPhrase(p.Localization[l].Description, sourcePhrases, translations),
		}
	}

	for i, p := range recipe.Preparation {
		recipe.Preparation[i].Localization[l] = types.PreparationLocalization{
			Title:        aiApplyPhrase(p.Localization[l].Title, sourcePhrases, translations),
			Instructions: aiApplyPhrase(p.Localization[l].Instructions, sourcePhrases, translations),
		}
		for j, g := range p.Ingredients {
			recipe.Preparation[i].Ingredients[j].Localization[l] = types.IngredientLocalization{
				Title: aiApplyPhrase(g.Localization[l].Title, sourcePhrases, translations),
			}
		}
	}

	return recipe
}

func aiApplyPhrase(orig string, sourcePhrases []string, translations []string) string {
	i := slices.Index(sourcePhrases, orig)
	if i > -1 {
		return translations[i]
	}
	return orig
}
