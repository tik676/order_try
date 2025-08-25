package router

import (
	"chat_service/internal/delivery/rest"
	websocket "chat_service/internal/delivery/webSocket"
	"chat_service/internal/middleware"
	"chat_service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter(usecase *usecase.UseCase, tokenManager middleware.TokenManager) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	httpRepo := rest.NewHTTPRepository(usecase)
	wsRepo := websocket.NewWebSocket(usecase)
	authMW := middleware.NewAuthMiddleware(tokenManager).OptionalAuth()

	r.GET("/messages", httpRepo.GetMessagesHandler)
	r.POST("/messages", authMW, httpRepo.PostMessageHandler)
	r.DELETE("/messages/:id", authMW, httpRepo.DeleteMessageHandler)
	r.GET("/ws", authMW, wsRepo.WsHandler)

	return r
}
