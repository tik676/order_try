package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
	"user_service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
	DB        *sql.DB
}

func NewJWTMaker(secretKey string, db *sql.DB) *JWTMaker {
	return &JWTMaker{secretKey: secretKey, DB: db}
}

func (j *JWTMaker) CreateToken(userID int64, role string) (*domain.Token, error) {
	now := time.Now()
	expires_at := now.Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     expires_at.Unix(),
		"iat":     now.Unix(),
		"iss":     "user-service",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accesstoken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Printf("Failed to create token:%v", err)
		return nil, err
	}

	query := `INSERT INTO tokens(access_token,refresh_token,created_at,expires_at)VALUES($1,$2,$3,$4)`
	_, err = j.DB.Exec(query, accesstoken, "", now, expires_at)
	if err != nil {
		log.Printf("Failed to save token:%v", err)
		return nil, err
	}

	return &domain.Token{
		AccessToken:  accesstoken,
		RefreshToken: "",
		CreatedAt:    now,
		ExpiresAt:    expires_at,
	}, nil

}

func (j *JWTMaker) VerifyToken(tokenString string) (userID int64, role string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = int64(claims["user_id"].(float64))
		role = claims["role"].(string)
		return userID, role, nil
	}

	return 0, "", errors.New("invalid token claims")
}

func (j *JWTMaker) RevokeToken(token string) error {
	query := `DELETE FROM tokens WHERE access_token=$1`
	_, err := j.DB.Exec(query, token)
	if err != nil {
		log.Printf("error token not found:%v", err)
		return err
	}
	return nil
}

func (j *JWTMaker) IsTokenValid(token string) (bool, error) {
	var tokesLife struct {
		createdAt time.Time
		expiresAt time.Time
	}
	query := `SELECT created_at,expires_at FROM tokens WHERE access_token = $1`
	err := j.DB.QueryRow(query, token).Scan(&tokesLife.createdAt, &tokesLife.expiresAt)
	if err != nil {
		log.Printf("error token not found:%v", err)
		return false, err
	}

	if time.Now().After(tokesLife.expiresAt) {
		return false, nil
	}
	if tokesLife.createdAt.After(time.Now()) {
		return false, nil
	}

	return true, nil
}
