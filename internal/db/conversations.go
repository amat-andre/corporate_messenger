package db

import (
	"fmt"

	"github.com/amat-andre/corporate_messenger/internal/models"
)

const getUserConversationsQuery = `
	SELECT
		c.id,
		c.type,
		c.created_at,
		u.id,
		u.login
	FROM conversations c
	INNER JOIN conversation_members current_member
		ON current_member.conversation_id = c.id
		AND current_member.user_id = $1
	INNER JOIN conversation_members other_member
		ON other_member.conversation_id = c.id
		AND other_member.user_id != $1
	INNER JOIN users u
		ON u.id = other_member.user_id
	WHERE c.type = 'direct'
	ORDER BY c.created_at DESC
	`

// GetUserConversations позволяет получить список диалогов в которых участвует пользователь
func (db *Database) GetUserConversations(userID int64) ([]*models.Conversation, error) {

	ctx, cancel := db.newContext()
	defer cancel()

	rows, err := db.DB.QueryContext(ctx, getUserConversationsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user conversations: %w", err)
	}
	defer rows.Close()

	conversations := make([]*models.Conversation, 0)

	for rows.Next() {

		conversation := &models.Conversation{
			OtherUser: &models.ConversationUser{},
		}

		err := rows.Scan(
			&conversation.ID,
			&conversation.Type,
			&conversation.CreatedAt,
			&conversation.OtherUser.ID,
			&conversation.OtherUser.Login,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}

		conversations = append(conversations, conversation)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate conversations: %w", err)
	}

	return conversations, nil
}
