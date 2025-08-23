package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenManager interface {
	VerifyToken(token string) (userID int64, role string, err error)
}

type AuthMiddleware struct {
	tokenManager TokenManager
}

func NewAuthMiddleware(tokenmanager TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenManager: tokenmanager}
}

func (am *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			userID, role, err := am.tokenManager.VerifyToken(tokenString)
			if err == nil {
				c.Set("user_id", userID)
				c.Set("role", role)
			}
		}
		c.Next()
	}
}
