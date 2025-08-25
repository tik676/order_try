package rest

import (
	"chat_service/internal/domain"
	"chat_service/internal/usecase"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HTTPRepository struct {
	usecase *usecase.UseCase
}

func NewHTTPRepository(usecase *usecase.UseCase) *HTTPRepository {
	return &HTTPRepository{usecase: usecase}
}

func (h *HTTPRepository) GetMessagesHandler(c *gin.Context) {
	limit := 50
	offset := 0

	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	messages, err := h.usecase.GetMessages(limit, offset)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get message"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *HTTPRepository) PostMessageHandler(c *gin.Context) {
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

	var req struct {
		Content string `json:"content"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse JSON"})
		return
	}

	msg, err := h.usecase.SendMessage(user, req.Content)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *HTTPRepository) DeleteMessageHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	if err := h.usecase.DeleteMessage(id); err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "succesfull"})
}
