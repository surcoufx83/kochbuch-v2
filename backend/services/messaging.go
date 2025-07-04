package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"kochbuch-v2-backend/types"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	connections      map[string]*wsConnection = make(map[string]*wsConnection)
	connectionsMutex sync.Mutex
	writeMutex       sync.Mutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins, adjust this for production
	},
}

type wsConnection struct {
	Connection       *websocket.Conn
	ConnectionParams AppParams
	User             *types.UserProfile
}

type wsIncomingAuthMessage struct {
	Token string `json:"token"`
}

type wsIncomingAuthCallbackMessage struct {
	State string `json:"state" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type wsIncomingIdMessage struct {
	Id   int       `json:"id" binding:"required"`
	Etag time.Time `json:"etag"`
}

type wsIncomingErrorReport struct {
	Url      string `json:"url" binding:"required"`
	Severity string `json:"severity" binding:"required"`
	Message  string `json:"error" binding:"required"`
}

type wsMessage struct {
	MsgType string           `json:"type" binding:"required"`
	Content string           `json:"content" binding:"required"`
	State   types.NullString `json:"state"`
}

type wsErrorContent struct {
	Code    int    `json:"error"`
	Message string `json:"message"`
}

type wsHelloContent struct {
	ConnectionParams AppParams          `json:"connection"`
	LoggedIn         bool               `json:"loggedIn"`
	User             *types.UserProfile `json:"user"`
}

type wsCategoriesContent struct {
	Categories map[uint16]types.Category `json:"categories"`
	Etag       string                    `json:"etag"`
}

type wsRecipesContent struct {
	Recipes map[uint32]*types.RecipeSimple `json:"recipes"`
	Etag    string                         `json:"etag"`
}

type wsUploadPictureContent struct {
	Recipe int32                        `json:"recipe" binding:"required"`
	Files  []wsUploadPictureFileContent `json:"files" binding:"required"`
}

type wsUploadPictureFileContent struct {
	Name          string `json:"name" binding:"required"`
	Type          string `json:"type" binding:"required"`
	Size          int    `json:"size" binding:"required"`
	Base64Content string `json:"base64" binding:"required"`
}

type wsUnitsContent struct {
	Units map[uint8]types.Unit `json:"units"`
	Etag  string               `json:"etag"`
}

func OnWebsocketConnect(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	// Wait for the first authentication message
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read failed:", err)
		return
	}

	// Decode the authentication message
	authMsg, err := wsReadMessage(msg)
	if err != nil {
		return
	}

	if authMsg.MsgType != "auth" || authMsg.Content == "" {
		log.Println("Invalid authentication message")
		return
	}

	var authToken wsIncomingAuthMessage
	if err := json.Unmarshal([]byte(authMsg.Content), &authToken); err != nil {
		log.Println("Error unmarshalling auth message:", err)
		return
	}

	state, params, _ := GetApplicationParams(conn, authToken.Token)
	code, _, user, err := GetSelfByState(state)
	if code != http.StatusUnauthorized && code != http.StatusOK {
		log.Println("Invalid token:", err)
	}

	connectionsMutex.Lock()
	connections[authToken.Token] = &wsConnection{
		Connection:       conn,
		ConnectionParams: params,
		User:             user,
	}
	connectionsMutex.Unlock()

	username := user.DisplayName
	if username == "" {
		username = "anonymous"
	}
	log.Printf("  > WS: %s connected", username)

	wsWelcome(connections[authToken.Token])

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read failed:", err)
			break
		}

		// log.Printf("  > WS: received msg %v", string(p))
		inMsg, err := wsReadMessage(p)
		if err != nil {
			break
		}
		go wsHandleMessage(connections[authToken.Token], inMsg)

	}
}

func OnUserProfileUpdated(state string, userProfile *types.UserProfile) {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()
	conn := connections[state]
	if conn.Connection == nil || conn.ConnectionParams.Session != state {
		return
	}

	conn.User = userProfile
	connections[state] = conn
	log.Printf("  > WS: set session %v to user %d", state, userProfile.Id)

	wsWelcome(connections[state])
}

func wsHandleMessage(conn *wsConnection, msg wsMessage) {
	log.Printf("  > WS: wsHandleMessage(%s)", msg.MsgType)
	switch msg.MsgType {

	case "bye":
		_, _ = Logout(conn)
		OnUserProfileUpdated(conn.ConnectionParams.Session, &types.UserProfile{})
		return

	case "categories_get_all":
		wsGetCategories(conn)
		return

	case "error_report":
		wsReportError(conn, msg)
		return

	case "oauth2_callback":
		wsOAuth2Callback(conn, msg)
		return

	case "recipe_get":
		wsGetRecipe(conn, msg)
		return

	case "recipes_get_all":
		wsGetRecipes(conn)
		return

	case "recipe_picture_upload":
		wsUploadPicture(conn, msg)
		return

	case "units_get_all":
		wsGetUnits(conn)
		return

	}
}

func wsGetCategories(conn *wsConnection) {
	categories, etag := GetCategories()
	var content wsCategoriesContent = wsCategoriesContent{
		Categories: categories,
		Etag:       etag,
	}
	jsoncontent, err := json.Marshal(content)
	if err != nil {
		log.Println("Error marshalling categories listing:", err)
		return
	}
	wsWriteMessage(conn, &wsMessage{
		MsgType: "categories",
		Content: string(jsoncontent),
	})
}

func wsGetRecipe(conn *wsConnection, msg wsMessage) {
	var data wsIncomingIdMessage
	if err := json.Unmarshal([]byte(msg.Content), &data); err != nil {
		log.Println("Error unmarshalling recipe request message:", err)
		wsWrite400BadRequest(conn, "recipe_get", types.NullString{Valid: false})
		return
	}

	recipe, err := GetRecipeWs(uint32(data.Id), conn)
	if err != nil {
		wsWrite403Forbidden(conn, "recipe_get", types.NullString{Valid: false})
		return
	}

	if data.Etag.Equal(recipe.ModifiedTime) {
		wsWrite304NotModified(conn, "recipe_get", types.NullString{Valid: false})
		return
	}

	jsoncontent, err := json.Marshal(recipe)
	if err != nil {
		log.Println("Error marshalling recipe response:", err)
		return
	}
	wsWriteMessage(conn, &wsMessage{
		MsgType: "recipe_get",
		Content: string(jsoncontent),
	})

}

func wsGetRecipes(conn *wsConnection) {
	recipes, etag := GetRecipes(conn.User)
	var content wsRecipesContent = wsRecipesContent{
		Recipes: recipes,
		Etag:    etag,
	}
	jsoncontent, err := json.Marshal(content)
	if err != nil {
		log.Println("Error marshalling recipes listing:", err)
		return
	}
	wsWriteMessage(conn, &wsMessage{
		MsgType: "recipes",
		Content: string(jsoncontent),
	})
}

func wsGetUnits(conn *wsConnection) {
	units, etag := GetUnits()
	var content wsUnitsContent = wsUnitsContent{
		Units: units,
		Etag:  etag,
	}
	jsoncontent, err := json.Marshal(content)
	if err != nil {
		log.Println("Error marshalling units listing:", err)
		return
	}
	wsWriteMessage(conn, &wsMessage{
		MsgType: "units",
		Content: string(jsoncontent),
	})
}

func wsNotifyRecipesChanged() {
	for _, conn := range connections {
		wsWriteMessage(conn, &wsMessage{
			MsgType: "recipes_etag",
			Content: recipesEtagStr,
		})
	}
}

func wsOAuth2Callback(conn *wsConnection, msg wsMessage) {
	var data wsIncomingAuthCallbackMessage
	if err := json.Unmarshal([]byte(msg.Content), &data); err != nil {
		wsWriteMessage(conn, &wsMessage{
			MsgType: "oauth2_response",
			Content: "400/Bad Request",
		})
		log.Println("Error unmarshalling callback message:", err)
		return
	}

	if data.State != conn.ConnectionParams.Session {
		wsWriteMessage(conn, &wsMessage{
			MsgType: "oauth2_response",
			Content: "400/Bad Request",
		})
		return
	}

	result, err := NcLoginCallback(data.State, data.Code)
	if err != nil {
		wsWriteMessage(conn, &wsMessage{
			MsgType: "oauth2_response",
			Content: "500/Internal Server Error",
		})
		return
	}

	if !result {
		wsWriteMessage(conn, &wsMessage{
			MsgType: "oauth2_response",
			Content: "403/Forbidden",
		})
		return
	}

	wsWriteMessage(conn, &wsMessage{
		MsgType: "oauth2_response",
		Content: "202/Accepted",
	})
}

func wsReadMessage(msg []byte) (wsMessage, error) {
	var message wsMessage
	if err := json.Unmarshal(msg, &message); err != nil {
		log.Println("Error unmarshalling message:", err)
		return wsMessage{}, err
	}
	return message, nil
}

func wsReportError(conn *wsConnection, msg wsMessage) {
	fn := "wsReportError()"
	var data wsIncomingErrorReport
	if err := json.Unmarshal([]byte(msg.Content), &data); err != nil {
		log.Printf("%v: Error unmarshalling recipe request message: %v", fn, err)
		wsWrite400BadRequest(conn, "error_report", types.NullString{Valid: false})
		return
	}

	if data.Severity != "I" && data.Severity != "E" && data.Severity != "W" {
		log.Printf("%v: Severity != I|E|W: %v", fn, data.Severity)
		data.Severity = "E"
	}

	stmt, err := dbPrepareStmt("wsReportError", "INSERT INTO `apilog`(`severity`, `reporter`, `host`, `agent`, `request_type`, `request_uri`, `request_length`, `message`) VALUES(?, 'Client', '', '', 'wss://', ?, 0, ?)")
	if err != nil {
		return
	}

	_, _ = stmt.Exec(data.Severity, data.Url, data.Message)
}

func wsUploadPicture(conn *wsConnection, msg wsMessage) {
	var data wsUploadPictureContent
	if err := json.Unmarshal([]byte(msg.Content), &data); err != nil {
		log.Println("Error unmarshalling request content:", err)
		wsWrite400BadRequest(conn, msg.MsgType, msg.State)
		return
	}

	fn := fmt.Sprintf("wsUploadPicture(%s, %d, %d files)", conn.User.DisplayName, data.Recipe, len(data.Files))

	recipe, err := GetRecipeWs(uint32(data.Recipe), conn)
	if err != nil {
		wsWrite403Forbidden(conn, msg.MsgType, msg.State)
		return
	}

	for i, file := range data.Files {
		// basic checks if all required properties exist
		if file.Name == "" || file.Size <= 0 || file.Base64Content == "" {
			log.Printf("%v: Invalid request data for i=%d", fn, i)
			wsWrite400BadRequest(conn, msg.MsgType, msg.State)
			return
		}

		// base64 decode of file content
		filecontent, err := base64.StdEncoding.DecodeString(file.Base64Content)
		if err != nil {
			log.Printf("%v: Base64 decode for i=%d failed with: %v", fn, i, err)
			wsWrite400BadRequest(conn, msg.MsgType, msg.State)
			return
		}

		// check base64 decoded string length vs. requested filesize
		if len(filecontent) != int(file.Size) {
			log.Printf("%v: Invalid filesize for i=%d: bytes=%d, request=%d", fn, i, len(filecontent), file.Size)
			wsWrite400BadRequest(conn, msg.MsgType, msg.State)
			return
		}

		hash := Sha256(fmt.Sprintf("%v,%v", string(filecontent), time.Now()), 8)
		subfolder := hash[:2]
		folder := fmt.Sprintf("/media/cbimages/%s", subfolder)
		filename := fmt.Sprintf("%s_%s", hash, file.Name)
		fullpath := fmt.Sprintf("/media/cbimages/%s/%s_%s", subfolder, hash, file.Name)
		pictureIndex := len(recipe.Pictures)

		log.Printf("%v: Fullpath for i=%d: %s", fn, i, fullpath)

		// recursive create image folder -> does nothing if exists already
		err = os.MkdirAll(folder, 0644)
		if err != nil {
			log.Printf("%v: Failed creating image folder for i=%d: %v", fn, i, err)
			wsWrite500InternalServerError(conn, msg.MsgType, msg.State)
			return
		}

		// sql transaction
		tx, err := Db.Begin()
		if err != nil {
			log.Printf("%v: Failed creating transaction: %v", fn, err)
			wsWrite500InternalServerError(conn, msg.MsgType, msg.State)
			_ = tx.Rollback()
			return
		}

		// preapre stmt
		stmt, err := dbPrepareStmt("wsUploadPicture_insert", "INSERT INTO `recipe_pictures`(`recipe_id`, `user_id`, `sortindex`, `name_de`, `name_en`, `name_fr`, `description_de`, `description_en`, `description_fr`, `hash`, `filename`, `fullpath`, `width`, `height`) VALUES(?, ?, ?, ?, ?, ?, '', '', '', ?, ?, ?, ?, ?)")
		if err != nil {
			log.Printf("%v: Failed preparing stmt: %v", fn, err)
			wsWrite500InternalServerError(conn, msg.MsgType, msg.State)
			_ = tx.Rollback()
			return
		}

		// assign stmt to transaction
		stmt = tx.Stmt(stmt)

		// save file content to disk
		err = os.WriteFile(fullpath, filecontent, 0644)
		if err != nil {
			log.Printf("%v: Failed saving file for i=%d: %v", fn, i, err)
			wsWrite500InternalServerError(conn, msg.MsgType, msg.State)
			_ = tx.Rollback()
			return
		}

		// open image to read dimensions
		imgFile, err := os.Open(fullpath)
		if err != nil {
			log.Printf("%v: Failed opening file for i=%d: %v", fn, i, err)
			wsWrite500InternalServerError(conn, msg.MsgType, msg.State)
			_ = tx.Rollback()
			_ = os.Remove(fullpath)
			return
		}
		imgFile.Seek(0, 0)

		// decode image to check dimensions
		img, _, err := image.Decode(imgFile)
		if err != nil {
			_ = imgFile.Close()
			log.Printf("%v: Failed decoding img file for i=%d: %v", fn, i, err)
			wsWrite500InternalServerError(conn, msg.MsgType, msg.State)
			_ = tx.Rollback()
			_ = os.Remove(fullpath)
			return
		}
		_ = imgFile.Close()

		// execute stmt
		result, err := stmt.Exec(recipe.Id, conn.User.Id, pictureIndex, file.Name, file.Name, file.Name, hash, filename, fullpath, img.Bounds().Dx(), img.Bounds().Dy())
		if err != nil {
			log.Printf("%v: Failed exec stmt: %v", fn, err)
			wsWrite500InternalServerError(conn, msg.MsgType, msg.State)
			_ = tx.Rollback()
			_ = os.Remove(fullpath)
			return
		}

		// retrieve picture id from database
		fileid, err := result.LastInsertId()
		if err != nil {
			log.Printf("%v: Failed retrieving insertId: %v", fn, err)
			wsWrite500InternalServerError(conn, msg.MsgType, msg.State)
			_ = tx.Rollback()
			_ = os.Remove(fullpath)
			return
		}

		// commit changes to the database
		err = tx.Commit()
		if err != nil {
			log.Printf("%v: Failed commiting stmt: %v", fn, err)
			wsWrite500InternalServerError(conn, msg.MsgType, msg.State)
			_ = tx.Rollback()
			_ = os.Remove(fullpath)
			return
		}

		log.Printf("%v: File i=%d created as id %d", fn, i, fileid)

		basename, ext := GetBasenameAndExtension(filename)
		picture := types.Picture{
			Id:       uint32(fileid),
			RecipeId: recipe.Id,
			User:     conn.User.SimpleProfile,
			Index:    uint8(pictureIndex),
			Localization: map[string]types.PictureLocalization{
				"de": {
					Name:        file.Name,
					Description: "",
				},
				"en": {
					Name:        file.Name,
					Description: "",
				},
				"fr": {
					Name:        file.Name,
					Description: "",
				},
			},
			FileName: filename,
			BaseName: basename,
			Ext:      ext,
			FullPath: fullpath,
			Uploaded: time.Now(),
			Dimension: types.PictureDimension{
				Height:         img.Bounds().Dy(),
				Width:          img.Bounds().Dx(),
				GeneratedSizes: []int{},
				Generated:      types.NullTime{Valid: false},
			},
		}
		AddPictureToRecipe(recipe, &picture)
	}

	wsWrite202Accepted(conn, msg.MsgType, msg.State)
	go touchRecipe(recipe)

}

func wsWelcome(conn *wsConnection) {

	var payload = wsHelloContent{
		ConnectionParams: conn.ConnectionParams,
		LoggedIn:         conn.User.Id > 0,
		User:             conn.User,
	}
	jsoncontent, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshalling recipes listing:", err)
		return
	}
	wsWriteMessage(conn, &wsMessage{
		MsgType: "hello",
		Content: string(jsoncontent),
	})

	wsWriteMessage(conn, &wsMessage{
		MsgType: "recipes_etag",
		Content: recipesEtagStr,
	})

	wsWriteMessage(conn, &wsMessage{
		MsgType: "categories_etag",
		Content: categoriesEtagStr,
	})

	wsWriteMessage(conn, &wsMessage{
		MsgType: "units_etag",
		Content: unitsEtagStr,
	})
}

func wsWelcomeAgain(user *types.UserProfile) {
	for _, conn := range connections {
		if conn.User != nil && conn.User.Id == user.Id {
			var payload = wsHelloContent{
				ConnectionParams: conn.ConnectionParams,
				LoggedIn:         conn.User.Id > 0,
				User:             conn.User,
			}
			jsoncontent, err := json.Marshal(payload)
			if err != nil {
				log.Println("Error marshalling recipes listing:", err)
				return
			}
			wsWriteMessage(conn, &wsMessage{
				MsgType: "hello",
				Content: string(jsoncontent),
			})
		}
	}
}

func wsWriteMessage(conn *wsConnection, message *wsMessage) {
	writeMutex.Lock()
	defer writeMutex.Unlock()
	msg, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("  > Failed json encode ws message: %v", message)
		return
	}
	if err := conn.Connection.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Println("Write failed:", err)
	}
}

func wsWrite202Accepted(conn *wsConnection, msgType string, state types.NullString) {
	var content wsErrorContent = wsErrorContent{
		Code:    202,
		Message: "Accepted",
	}
	jsoncontent, err := json.Marshal(content)
	if err != nil {
		log.Println("Error marshalling recipe response:", err)
		return
	}
	wsWriteMessage(conn, &wsMessage{
		MsgType: msgType,
		Content: string(jsoncontent),
		State:   state,
	})
}

func wsWrite304NotModified(conn *wsConnection, msgType string, state types.NullString) {
	var content wsErrorContent = wsErrorContent{
		Code:    304,
		Message: "NotModified",
	}
	jsoncontent, err := json.Marshal(content)
	if err != nil {
		log.Println("Error marshalling recipe response:", err)
		return
	}
	wsWriteMessage(conn, &wsMessage{
		MsgType: msgType,
		Content: string(jsoncontent),
		State:   state,
	})
}

func wsWrite400BadRequest(conn *wsConnection, msgType string, state types.NullString) {
	var content wsErrorContent = wsErrorContent{
		Code:    400,
		Message: "Bad request",
	}
	jsoncontent, err := json.Marshal(content)
	if err != nil {
		log.Println("Error marshalling recipe response:", err)
		return
	}
	wsWriteMessage(conn, &wsMessage{
		MsgType: msgType,
		Content: string(jsoncontent),
		State:   state,
	})
}

func wsWrite403Forbidden(conn *wsConnection, msgType string, state types.NullString) {
	var content wsErrorContent = wsErrorContent{
		Code:    403,
		Message: "Forbidden",
	}
	jsoncontent, err := json.Marshal(content)
	if err != nil {
		log.Println("Error marshalling recipe response:", err)
		return
	}
	wsWriteMessage(conn, &wsMessage{
		MsgType: msgType,
		Content: string(jsoncontent),
		State:   state,
	})
}

func wsWrite404NotFound(conn *wsConnection, msgType string, state types.NullString) {
	var content wsErrorContent = wsErrorContent{
		Code:    404,
		Message: "Not Found",
	}
	jsoncontent, err := json.Marshal(content)
	if err != nil {
		log.Println("Error marshalling recipe response:", err)
		return
	}
	wsWriteMessage(conn, &wsMessage{
		MsgType: msgType,
		Content: string(jsoncontent),
		State:   state,
	})
}

func wsWrite500InternalServerError(conn *wsConnection, msgType string, state types.NullString) {
	var content wsErrorContent = wsErrorContent{
		Code:    500,
		Message: "Internal Server Error",
	}
	jsoncontent, err := json.Marshal(content)
	if err != nil {
		log.Println("Error marshalling recipe response:", err)
		return
	}
	wsWriteMessage(conn, &wsMessage{
		MsgType: msgType,
		Content: string(jsoncontent),
		State:   state,
	})
}
