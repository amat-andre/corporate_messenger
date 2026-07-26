package db

import (
	"context"
	"database/sql"
	"fmt"
)

func (db *Database) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	err = fn(tx)
	if err != nil {

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("rollback failed: %v: %w", rollbackErr, err)
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
