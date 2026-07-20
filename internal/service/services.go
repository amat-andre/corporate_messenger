package service

import (
	"log/slog"

	"github.com/amat-andre/corporate_messenger/internal/db"
)

type Services struct {
	Auth *AuthService
	Chat *ChatService
	User *UserService
}

func NewServices(db *db.Database, logger *slog.Logger) *Services {
	return &Services{
		Auth: NewAuthService(db, logger),
		User: NewUserService(db, logger),
		Chat: NewChatService(db, logger),
	}
}