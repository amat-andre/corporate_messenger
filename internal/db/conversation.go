package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/amat-andre/corporate_messenger/internal/models"
)

var ErrConversationNotFound = errors.New("conversation not found")

// CreateConversation создает диалог
func (db *Database) CreateConversation(conversation *models.Conversation) error {
	const query = `
		INSERT INTO conversations (type)
		VALUES ($1)
		RETURNING id, created_at
	`

	ctx, cancel := db.newContext()
	defer cancel()

	err := db.DB.QueryRowContext(
		ctx,
		query,
		conversation.Type,
	).Scan(&conversation.ID, &conversation.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	return nil
}

// AddConversationMember добавляет участника в диалог
func (db *Database) AddConversationMember(member *models.ConversationMember) error {
	const query = `
		INSERT INTO conversation_members (conversation_id, user_id)
		VALUES ($1, $2)
		RETURNING joined_at
	`

	ctx, cancel := db.newContext()
	defer cancel()

	err := db.DB.QueryRowContext(
		ctx,
		query,
		member.ConversationID,
		member.UserID,
	).Scan(&member.JoinedAt)

	if err != nil {
		return fmt.Errorf("failed to add conversation member: %w", err)
	}

	return nil
}

// ConversationExists проверяет существование диалога
func (db *Database) ConversationExists(conversationID int64) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM conversations
			WHERE id = $1
		)
	`

	ctx, cancel := db.newContext()
	defer cancel()

	var exists bool

	err := db.DB.QueryRowContext(
		ctx,
		query,
		conversationID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check conversation: %w", err)
	}

	return exists, nil
}

// IsConversationMember проверяет состоит ли пользователь в диалоге
func (db *Database) IsConversationMember(conversationID, userID int64) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM conversation_members
			WHERE conversation_id = $1 AND user_id = $2
		)
	`

	ctx, cancel := db.newContext()
	defer cancel()

	var exists bool

	err := db.DB.QueryRowContext(
		ctx,
		query,
		conversationID,
		userID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check conversation member: %w", err)
	}

	return exists, nil
}

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

		if err == sql.ErrNoRows {
			return nil, ErrConversationNotFound
		}

		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return conversation, nil
}

// CreateDirectConversation создает личный диалог между двумя пользователями
func (db *Database) CreateDirectConversation(firstUserID, secondUserID int64) (*models.Conversation, error) {

	ctx, cancel := db.newContext()
	defer cancel()

	var conversation *models.Conversation

	err := db.WithTransaction(ctx, func(tx *sql.Tx) error {

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
