package middleware

import (
	"strings"
	"user_service/internal/domain"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenManager domain.TokenManager
}

func NewAuthMiddleware(tokenmanager domain.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenManager: tokenmanager}
}

func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(401, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		userID, name, role, err := am.tokenManager.VerifyToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("name", name)
		c.Set("role", role)

		c.Next()
	}
}
