package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/amat-andre/corporate_messenger/internal/models"
)

func (db *Database) CreateMessage(conversationID, userID int64, text string) (*models.Message, error) {

	ctx, cancel := db.newContext()
	defer cancel()

	const query = `
		WITH inserted_message AS (
			INSERT INTO messages (
				user_id,
				conversation_id,
				text
			)
			SELECT
				$1,
				$2,
				$3
			WHERE EXISTS (
				SELECT 1
				FROM conversation_members
				WHERE conversation_id = $2
					AND user_id = $1
			)
			RETURNING
				id,
				user_id,
				conversation_id,
				text,
				created_at
		)
		SELECT
			m.id,
			m.user_id,
			m.conversation_id,
			u.login,
			m.text,
			m.created_at
		FROM inserted_message m
		INNER JOIN users u
			ON u.id = m.user_id
	`

	message := &models.Message{}

	err := db.DB.QueryRowContext(
		ctx,
		query,
		userID,
		conversationID,
		text,
	).Scan(&message.ID, &message.UserID, &message.ConversationID, &message.Username, &message.Text, &message.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrConversationNotFound
		}
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return message, nil
}

func (db *Database) GetConversationMessages(conversationID, userID int64) ([]*models.Message, error) {

	ctx, cancel := db.newContext()
	defer cancel()

	const query = `
		SELECT
			m.id,
			m.user_id,
			m.conversation_id,
			u.login,
			m.text,
			m.created_at
		FROM messages m
		INNER JOIN users u
			ON u.id = m.user_id
		INNER JOIN conversation_members cm
			ON cm.conversation_id = m.conversation_id
			AND cm.user_id = $2
		WHERE m.conversation_id = $1
		ORDER BY m.created_at ASC, m.id ASC
	`

	rows, err := db.DB.QueryContext(ctx, query, conversationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation messages: %w", err)
	}
	defer rows.Close()

	messages := make([]*models.Message, 0)

	for rows.Next() {

		message := &models.Message{}

		err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.ConversationID,
			&message.Username,
			&message.Text,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation message: %w", err)
		}

		messages = append(messages, message)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate conversation messages: %w", err)
	}

	if len(messages) == 0 {
		var exists bool

		err := db.DB.QueryRowContext(ctx,
			`
				SELECT EXISTS (
					SELECT 1
					FROM conversation_members
					WHERE conversation_id = $1
						AND user_id = $2
				)
			`,
			conversationID,
			userID,
		).Scan(&exists)

		if err != nil {
			return nil, fmt.Errorf("failed to check conversation membership: %w", err)
		}

		if !exists {
			return nil, ErrConversationNotFound
		}
	}

	return messages, nil
}
