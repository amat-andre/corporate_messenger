package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/amat-andre/corporate_messenger/internal/db"
	"github.com/amat-andre/corporate_messenger/internal/models"
)

const maxMessageTextLength = 5000

type MessageService struct {
	db     *db.Database
	logger *slog.Logger
}

func NewMessageService(db *db.Database, logger *slog.Logger) *MessageService {
	return &MessageService{
		db:     db,
		logger: logger,
	}
}

func (s *MessageService) CreateMessage(conversationID, userID int64, text string) (*models.Message, error) {

	text = strings.TrimSpace(text)

	if text == "" {
		return nil, ErrMessageTextEmpty
	}

	if len([]rune(text)) > maxMessageTextLength {
		return nil, ErrMessageTextTooLong
	}

	message, err := s.db.CreateMessage(conversationID, userID, text)
	if err != nil {
		if errors.Is(err, db.ErrConversationNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	s.logger.Info("message created",
		"message_id", message.ID,
		"conversation_id", message.ConversationID,
		"user_id", message.UserID,
	)

	return message, nil
}

func (s *MessageService) GetConversationMessages(conversationID, userID int64) ([]*models.Message, error) {

	messages, err := s.db.GetConversationMessages(conversationID, userID)
	if err != nil {
		if errors.Is(err, db.ErrConversationNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, fmt.Errorf("failed to get conversation messages: %w", err)
	}

	return messages, nil
}
