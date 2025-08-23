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
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to upgrade")
		return
	}
	defer conn.Close()

	clientManager.AddClient(conn)
	defer func() {
		clientManager.RemoveClient(conn)
		conn.Close()
	}()
	var user domain.User

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.String(http.StatusBadRequest, "Failed user not exists")
		return
	}
	userID, ok := userIDValue.(int64)
	if !ok {
		c.String(http.StatusBadRequest, "Failed to upgrade")
	}

	roleValue, ok := c.Get("role")
	if !ok {
		c.String(http.StatusBadRequest, "Failed to get a role")
		return
	}
	role, ok := roleValue.(string)

	user.ID = userID
	user.Role = role

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		msg, err := ws.usecase.SendMessage(user, string(message))
		if err != nil {
			log.Printf("Failed to send message for user %d: %v", user.ID, err)
			errMsg := map[string]string{"error": "Failed to send message"}
			errJSON, _ := json.Marshal(errMsg)
			conn.WriteMessage(websocket.TextMessage, errJSON)
			continue
		}

		msgJSON, err := json.Marshal(msg)
		if err != nil {
			log.Println("Failed to marshal message:", err)
			continue
		}

		clientManager.Broadcast(msgJSON)
	}
}
