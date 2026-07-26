package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/amat-andre/corporate_messenger/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

// CreateUser сохраняет нового пользователя.
func (db *Database) CreateUser(user *models.User) error {
	const query = `
		INSERT INTO users (login, password_hash, created_at)
		VALUES ($1, $2, $3)
		RETURNING id 
	`
	// возможно убрать RETURNING если не понадобится
	ctx, cancel := db.newContext()
	defer cancel()

	err := db.DB.QueryRowContext(
		ctx,
		query,
		user.Login,
		user.PasswordHash,
		user.CreatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByLogin ищет пользователя по логину.
func (db *Database) GetUserByLogin(login string) (*models.User, error) {
	const query = `
		SELECT id, login, password_hash, created_at
		FROM users
		WHERE login = $1
	`

	ctx, cancel := db.newContext()
	defer cancel()

	user := &models.User{}

	err := db.DB.QueryRowContext(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID ищет пользователя по ID.
func (db *Database) GetUserByID(id int64) (*models.User, error) {
	const query = `
		SELECT id, login, password_hash, created_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := db.newContext()
	defer cancel()

	user := &models.User{}

	err := db.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
