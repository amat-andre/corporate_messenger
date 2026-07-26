package models

import "time"

// Conversation представляет личный чат между пользователями.
type Conversation struct {
	ID        int64 `db:"conversation_id" json:"conversation_id"`
	Type      string
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	OtherUser *ConversationUser
}

// ConversationUser - безопасная модель пользователя.
type ConversationUser struct {
	ID    int64
	Login string
}

// ConversationMember представляет участие пользователя в чате.
type ConversationMember struct {
	ConversationID int64 `db:"conversation_id" json:"conversation_id"`
	UserID         int64 `db:"user_id" json:"user_id"`
	JoinedAt       time.Time
}
