package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenManager interface {
	VerifyToken(token string) (userID int64, name, role string, err error)
}

type AuthMiddleware struct {
	tokenManager TokenManager
}

func NewAuthMiddleware(tokenmanager TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenManager: tokenmanager}
}

func (am *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := int64(0)
		name := "Аноним"
		role := "anonym"
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if verifiedUserID, verifiedName, verifiedRole, err := am.tokenManager.VerifyToken(tokenString); err == nil {
				userID = verifiedUserID
				name = verifiedName
				role = verifiedRole
			}
		}
		c.Set("user_id", userID)
		c.Set("name", name)
		c.Set("role", role)
		c.Next()
	}
}
