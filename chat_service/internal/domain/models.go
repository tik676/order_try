package domain

import "time"

type User struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	IsAnon bool   `json:"is_anon"`
}

type Message struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Content    string    `json:"content"`
	Created_at time.Time `json:"created_at"`
}

type MessageRepository interface {
	SendMessage(msg Message) (Message, error)
	DeleteMessage(id int64) error
	GetMessages(limit, offset int) ([]Message, error)
}
