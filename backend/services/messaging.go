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
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins, adjust this for production
	},
}

type wsConnection struct {
	Connection *websocket.Conn
	User       types.UserProfile
}

type wsIncomingAuthMessage struct {
	Token string `json:"token"`
}

type wsMessage struct {
	MsgType string `json:"type"`
	Content string `json:"content"`
}

type wsRecipesContent struct {
	Recipes map[uint32]*types.RecipeSimple `json:"recipes"`
	Etag    string                         `json:"etag"`
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

	code, _, user, err := GetSelfByState(authToken.Token)
	if code != http.StatusUnauthorized && code != http.StatusOK {
		log.Println("Invalid token:", err)
		return
	}

	connectionsMutex.Lock()
	connections[authToken.Token] = &wsConnection{
		Connection: conn,
		User:       user,
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

		/* if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("Write failed:", err)
			break
		} */
	}
}

func wsHandleMessage(conn *wsConnection, msg wsMessage) {
	log.Println("  > WS: async handle message")
	switch msg.MsgType {
	case "recipes_get_all":
		wsGetRecipes(conn)
		return
	}
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

func wsReadMessage(msg []byte) (wsMessage, error) {
	var message wsMessage
	if err := json.Unmarshal(msg, &message); err != nil {
		log.Println("Error unmarshalling message:", err)
		return wsMessage{}, err
	}
	return message, nil
}

func wsWelcome(conn *wsConnection) {
	message := fmt.Sprintf("Hello %s", conn.User.DisplayName)
	wsWriteMessage(conn, &wsMessage{
		MsgType: "none",
		Content: message,
	})
	wsWriteMessage(conn, &wsMessage{
		MsgType: "recipes_etag",
		Content: recipesEtagStr,
	})
	wsWriteMessage(conn, &wsMessage{
		MsgType: "categories_etag",
		Content: categoriesEtagStr,
	})
}

func wsWriteMessage(conn *wsConnection, message *wsMessage) {
	msg, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("  > Failed json encode ws message: %v", message)
		return
	}
	if err := conn.Connection.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Println("Write failed:", err)
	}
}
