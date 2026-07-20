package models

import "time"

// User представляет пользователя системы.
type User struct {
	ID       	 int64  	`db:"id" json:"id"`
	Login   	 string 	`db:"login" json:"login"`
	PasswordHash string 	`db:"password_hash" json:"-"`
	CreatedAt 	 time.Time 	`db:"created_at" json:"created_at"`
}

// Message представляет сообщение в чате.
type Message struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Username  string    `db:"username" json:"username"`
	Text      string    `db:"text" json:"text"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Conversation представляет чат между пользователями.
type Conversation struct {
	ID        int64		`db:"conversation_id" json:"conversation_id"`
	CreatedAt time.Time	`db:"created_at" json:"created_at"`
}

// ConversationMember представляет связь между пользователем и чатом.
type ConversationMember struct {
	ConversationID int64	`db:"conversation_id" json:"conversation_id"`
	UserID         int64	`db:"user_id" json:"user_id"`
}

func NewUser(login string, passwordHash string) *User {
	return &User{
		Login:        login,
		PasswordHash: passwordHash, // пока оставляем, посмотреть по логике
		CreatedAt:    time.Now(),
	}
}