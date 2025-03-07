package services

import (
	"encoding/json"
	"fmt"
	"kochbuch-v2-backend/types"
	"log"
	"net/http"
	"sync"

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
	User             types.UserProfile
}

type wsIncomingAuthMessage struct {
	Token string `json:"token"`
}
type wsIncomingAuthCallbackMessage struct {
	State string `json:"state" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type wsMessage struct {
	MsgType string `json:"type"`
	Content string `json:"content"`
}

type wsHelloContent struct {
	ConnectionParams AppParams         `json:"connection"`
	LoggedIn         bool              `json:"loggedIn"`
	User             types.UserProfile `json:"user"`
}

type wsCategoriesContent struct {
	Categories map[uint16]types.Category `json:"categories"`
	Etag       string                    `json:"etag"`
}

type wsRecipesContent struct {
	Recipes map[uint32]*types.RecipeSimple `json:"recipes"`
	Etag    string                         `json:"etag"`
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

	log.Println(authMsg.Content)

	var authToken wsIncomingAuthMessage
	if err := json.Unmarshal([]byte(authMsg.Content), &authToken); err != nil {
		log.Println("Error unmarshalling auth message:", err)
		return
	}

	log.Println(authToken.Token)

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

		log.Printf("  > WS: received msg %v", string(p))
		inMsg, err := wsReadMessage(p)
		if err != nil {
			break
		}
		go wsHandleMessage(connections[authToken.Token], inMsg)

	}
}

func OnUserProfileUpdated(state string, userProfile types.UserProfile) {
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
	log.Println("  > WS: async handle message")
	switch msg.MsgType {

	case "bye":
		_, _ = Logout(conn)
		OnUserProfileUpdated(conn.ConnectionParams.Session, types.UserProfile{})
		return

	case "categories_get_all":
		wsGetCategories(conn)
		return

	case "oauth2_callback":
		wsOAuth2Callback(conn, msg)
		return

	case "recipes_get_all":
		wsGetRecipes(conn)
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

func wsGetRecipes(conn *wsConnection) {
	recipes, etag := GetRecipes(&conn.User)
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
