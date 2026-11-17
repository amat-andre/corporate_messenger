package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/amat-andre/corporate_messenger/internal/models"
)

var ErrConversationNotFound = errors.New("conversation not found")

// GetConversationByID получает диалог по ID
func (db *Database) GetConversationByID(conversationID, userID int64) (*models.Conversation, error) {
	const query = `SELECT
			c.id,
			c.type,
			c.created_at,
			u.id,
			u.login
		FROM conversations c
		INNER JOIN conversation_members current_member
			ON current_member.conversation_id = c.id
			AND current_member.user_id = $2
		INNER JOIN conversation_members other_member
			ON other_member.conversation_id = c.id
			AND other_member.user_id != $2
		INNER JOIN users u
			ON u.id = other_member.user_id
		WHERE c.id = $1
			AND c.type = 'direct'
	`

	ctx, cancel := db.newContext()
	defer cancel()

	conversation := &models.Conversation{
		OtherUser: &models.ConversationUser{},
	}

	err := db.DB.QueryRowContext(
		ctx,
		query,
		conversationID,
		userID,
	).Scan(
		&conversation.ID,
		&conversation.Type,
		&conversation.CreatedAt,
		&conversation.OtherUser.ID,
		&conversation.OtherUser.Login,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrConversationNotFound
		}

		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return conversation, nil
}

// CreateDirectConversation создает личный диалог между двумя пользователями.
// Для пары пользователей берется advisory lock, поэтому параллельные
// вызовы не создают дублирующихся диалогов.
func (db *Database) CreateDirectConversation(firstUserID, secondUserID int64) (*models.Conversation, error) {

	ctx, cancel := db.newContext()
	defer cancel()

	var conversation *models.Conversation

	err := db.WithTransaction(ctx, func(tx *sql.Tx) error {

		// Блокируем пару пользователей на время транзакции,
		// чтобы find-or-create выполнялся строго последовательно.
		lockKey := fmt.Sprintf("%d:%d", min(firstUserID, secondUserID), max(firstUserID, secondUserID))

		const lockQuery = `SELECT pg_advisory_xact_lock(hashtext($1))`

		if _, err := tx.ExecContext(ctx, lockQuery, lockKey); err != nil {
			return fmt.Errorf("failed to acquire conversation lock: %w", err)
		}

		conversation = &models.Conversation{}

		const findDirectConversationQuery = `
		SELECT
			c.id,
			c.type,
			c.created_at
		FROM conversations c
		JOIN conversation_members cm1
			ON cm1.conversation_id = c.id
		JOIN conversation_members cm2
			ON cm2.conversation_id = c.id
		WHERE c.type = 'direct'
			AND cm1.user_id = $1
			AND cm2.user_id = $2
	`

		err := tx.QueryRowContext(
			ctx,
			findDirectConversationQuery,
			firstUserID,
			secondUserID,
		).Scan(&conversation.ID, &conversation.Type, &conversation.CreatedAt)

		if err == nil {
			return nil
		}

		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to find direct conversation: %w", err)
		}

		conversation.Type = "direct"

		const createConversationQuery = `
			INSERT INTO conversations (type)
			VALUES ($1)
			RETURNING id, created_at
		`

		err = tx.QueryRowContext(
			ctx,
			createConversationQuery,
			conversation.Type,
		).Scan(&conversation.ID, &conversation.CreatedAt)

		if err != nil {
			return fmt.Errorf("failed to create conversation: %w", err)
		}

		const addConversationMemberQuery = `
			INSERT INTO conversation_members (conversation_id, user_id)
			VALUES ($1, $2)
		`

		_, err = tx.ExecContext(ctx, addConversationMemberQuery, conversation.ID, firstUserID)
		if err != nil {
			return fmt.Errorf("failed to add first conversation member: %w", err)
		}

		_, err = tx.ExecContext(ctx, addConversationMemberQuery, conversation.ID, secondUserID)
		if err != nil {
			return fmt.Errorf("failed to add second conversation member: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return conversation, nil
}
