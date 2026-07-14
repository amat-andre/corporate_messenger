package service

import (
	"log/slog"

	"github.com/amat-andre/corporate_messenger/internal/db"
)

type ChatService struct {
	db     *db.Database
	logger *slog.Logger
}

func NewChatService(db *db.Database, logger *slog.Logger) *ChatService {
	return &ChatService{
		db:     db,
		logger: logger,
	}
}

func (s *ChatService) CreateConversation(ownerID int64, users []int64) (int64, error){
	return 0, nil
}