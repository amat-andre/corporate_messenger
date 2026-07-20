package db

import (
	"fmt"
)

const schema = `
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_messages_created_at
ON messages(created_at);

CREATE INDEX IF NOT EXISTS idx_messages_user_id
ON messages(user_id);

CREATE TABLE IF NOT EXISTS conversations (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS conversation_members (
    conversation_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,

    PRIMARY KEY(conversation_id, user_id),

    FOREIGN KEY(conversation_id)
        REFERENCES conversations(id)
        ON DELETE CASCADE,

    FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
`

func (base *Database) Migrate() error {
	ctx, cancel := base.newContext()
	defer cancel()

	_, err := base.DB.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to apply schema: %w", err)
	}

	base.logger.Info("database schema applied")

	return nil
}