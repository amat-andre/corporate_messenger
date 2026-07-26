package dto

import "time"

type CreateMessageRequest struct {
	Text string `json:"text"`
}

// Структура сообщения для ответа
type MessageResponse struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	ConversationID int64     `json:"conversation_id"`
	Username       string    `json:"username"`
	Text           string    `json:"text"`
	CreatedAt      time.Time `json:"created_at"`
}

// Структура для списка сообщений
type MessagesResponse struct {
	Messages []MessageResponse `json:"messages"`
}
