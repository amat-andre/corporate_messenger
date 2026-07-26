package service

import (
	"log/slog"

	"github.com/amat-andre/corporate_messenger/internal/db"
)

type Services struct {
	Auth         *AuthService
	User         *UserService
	Conversation *ConversationService
	Message      *MessageService
}

func NewServices(db *db.Database, logger *slog.Logger) *Services {
	return &Services{
		Auth:         NewAuthService(db, logger),
		User:         NewUserService(db, logger),
		Conversation: NewConversationService(db, logger),
		Message:      NewMessageService(db, logger),
	}
}
