package dto

import "time"

// Структура для получения ID пользователя с которым необходимо начать личный диалог
type CreateDirectConversationRequest struct {
	UserID int64 `json:"user_id"`
}

// Структура собеседника
type ConversationUserResponse struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

// Структура диалога
type ConversationResponse struct {
	ID        int64                     `json:"id"`
	Type      string                    `json:"type"`
	CreatedAt time.Time                 `json:"created_at"`
	User      *ConversationUserResponse `json:"user,omitempty"`
}

// Структура для получения пользователем списка диалогов в которых он состоит
type ConversationsResponse struct {
	Conversations []ConversationResponse `json:"conversations"`
}
