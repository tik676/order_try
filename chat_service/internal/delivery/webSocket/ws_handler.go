package websocket

import (
	"chat_service/internal/domain"
	"chat_service/internal/usecase"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketRepository struct {
	usecase *usecase.UseCase
}

func NewWebSocket(usecase *usecase.UseCase) *WebSocketRepository {
	return &WebSocketRepository{usecase: usecase}
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clientManager = NewClientManager()
)

func (ws *WebSocketRepository) WsHandler(c *gin.Context) {
	var user domain.User

	userIDRaw, _ := c.Get("user_id")
	nameRaw, _ := c.Get("name")
	roleRaw, _ := c.Get("role")

	userID, ok := userIDRaw.(int64)
	if !ok {
		userID = 0
	}

	name, ok := nameRaw.(string)
	if !ok {
		name = "Аноним"
	}

	role, ok := roleRaw.(string)
	if !ok {
		role = "anonym"
	}

	user.ID = userID
	user.Name = name
	user.Role = role
	user.IsAnon = (userID == 0)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to upgrade")
		return
	}
	defer conn.Close()

	clientManager.AddClient(conn)
	defer clientManager.RemoveClient(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		msg, err := ws.usecase.SendMessage(user, string(message))
		if err != nil {
			log.Printf("Error saving message: %v", err)
			errMsg := map[string]string{"error": "Failed to send message"}
			errJSON, _ := json.Marshal(errMsg)
			conn.WriteMessage(websocket.TextMessage, errJSON)
			continue
		}

		log.Printf("Message saved: ID=%d, UserID=%d, Content=%s", msg.ID, msg.UserID, msg.Content)

		msgJSON, err := json.Marshal(msg)
		if err != nil {
			log.Println("Failed to marshal message:", err)
			continue
		}
		clientManager.Broadcast(msgJSON)
	}
}
