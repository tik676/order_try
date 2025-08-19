package rest

import (
	"user_service/internal/domain"
	"user_service/internal/middleware"
	"user_service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter(usecase *usecase.UseCase, tokenManager domain.TokenManager) *gin.Engine {
	r := gin.Default()
	handler := NewHTTPHandler(usecase)
	middleware := middleware.NewAuthMiddleware(tokenManager)

	r.POST("/register", handler.RegisterHandler)
	r.POST("/login", handler.LoginHandler)
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth())
	{

	}

	return r
}
