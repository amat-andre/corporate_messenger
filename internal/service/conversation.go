package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/amat-andre/corporate_messenger/internal/db"
	"github.com/amat-andre/corporate_messenger/internal/models"
)

type ConversationService struct {
	db     *db.Database
	logger *slog.Logger
}

func NewConversationService(db *db.Database, logger *slog.Logger) *ConversationService {
	return &ConversationService{
		db:     db,
		logger: logger,
	}
}

// CreateDirectConversation создает диалог между двумя пользователями
func (s *ConversationService) CreateDirectConversation(firstUserID, secondUserID int64) (*models.Conversation, error) {

	if firstUserID == secondUserID {
		return nil, ErrConversationWithYourself
	}

	_, err := s.db.GetUserByID(secondUserID)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find conversation user: %w", err)
	}

	conversation, err := s.db.CreateDirectConversation(firstUserID, secondUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create direct conversation: %w", err)
	}

	s.logger.Info("direct conversation created", "conversation_id", conversation.ID,
		"first_user_id", firstUserID,
		"second_user_id", secondUserID)

	return conversation, nil
}

// GetUserConversations получает список диалогов пользователя
func (s *ConversationService) GetUserConversations(userID int64) ([]*models.Conversation, error) {

	conversations, err := s.db.GetUserConversations(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user conversations: %w", err)
	}

	return conversations, nil
}

// GetConversationByID получает диалог по ID. Пользователь может получить только диалог в котором он состоит.
func (s *ConversationService) GetConversationByID(conversationID, userID int64) (*models.Conversation, error) {

	conversation, err := s.db.GetConversationByID(conversationID, userID)
	if err != nil {
		if errors.Is(err, db.ErrConversationNotFound) {
			return nil, ErrConversationNotFound
		}

		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return conversation, nil
}
