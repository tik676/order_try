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

	httpRepo := rest.NewHTTPRepository(usecase)
	wsRepo := websocket.NewWebSocket(usecase)
	authMW := middleware.NewAuthMiddleware(tokenManager).OptionalAuth()

	r.GET("/messages", httpRepo.GetMessagesHandler)
	r.POST("/messages", authMW, httpRepo.PostMessageHandler)
	r.DELETE("/messages/:id", authMW, httpRepo.DeleteMessageHandler)
	r.GET("/ws", authMW, wsRepo.WsHandler)

	return r
}
