package middleware

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secret string) *JWTMaker {
	return &JWTMaker{secretKey: secret}
}

func (j *JWTMaker) VerifyToken(tokenString string) (userID int64, name, role string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return 0, "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat, ok1 := claims["user_id"].(float64)
		role, ok2 := claims["role"].(string)
		name, ok3 := claims["name"].(string)
		expFloat, ok4 := claims["exp"].(float64)
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return 0, "", "", errors.New("invalid token claims")
		}
		if time.Now().Unix() > int64(expFloat) {
			return 0, "", "", errors.New("token expired")
		}
		return int64(userIDFloat), name, role, nil
	}

	return 0, "", "", errors.New("invalid token")
}
