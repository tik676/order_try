package domain

import "time"

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"password_Hash"`
	Role         string    `json:"role"`
	RegisteredAt time.Time `json:"registered_at"`
}

type AuthorizationInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type Authorization interface {
	Register(input AuthorizationInput) (*User, error)
	Login(input AuthorizationInput) (*User, error)
}

type TokenManager interface {
	CreateToken(userID int64, role string) (*Token, error)
	VerifyToken(token string) (userID int64, role string, err error)
	RevokeToken(token string) error
	IsTokenValid(token string) (bool, error)
}
