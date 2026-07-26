package db

import (
	"fmt"

	"github.com/amat-andre/corporate_messenger/internal/models"
)

func (db *Database) SearchUsers(search string) ([]*models.User, error) {
	const query = `
		SELECT id, login, created_at
		FROM users
		WHERE login ILIKE '%' || $1 || '%'
		ORDER BY login
		LIMIT 20
	`

	ctx, cancel := db.newContext()
	defer cancel()

	rows, err := db.DB.QueryContext(ctx, query, search)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	users := []*models.User{}

	for rows.Next() {
		user := &models.User{}

		err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate users: %w", err)
	}

	return users, nil
}
