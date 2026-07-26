package models

import "time"

// Message представляет сообщение в чате.
type Message struct {
	ID             int64     `db:"id" json:"id"`
	UserID         int64     `db:"user_id" json:"user_id"`
	ConversationID int64     `db:"conversation_id" json:"conversation_id"`
	Username       string    `db:"username" json:"username"`
	Text           string    `db:"text" json:"text"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}
