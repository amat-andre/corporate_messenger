package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"log/slog"

	"github.com/amat-andre/corporate_messenger/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	DB *sql.DB
	logger *slog.Logger
}

func New(logger *slog.Logger) *Database {
	return &Database{
		logger: logger,
	}
}

func (base *Database) newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func (base *Database) Connect(cfg config.DBConfig) error {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := base.newContext()
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	base.DB = db
	base.logger.Info("database connected")

	return nil
}

func (base *Database) Close() {
	if base.DB == nil {
		return
	}

	if err := base.DB.Close(); err != nil {
		base.logger.Error("failed to close database:", "error", err)
	} else {
		base.logger.Info("database closed")
	}
}

func (base *Database) CreateConversation(userIDs []int64) (int64, error) {
	return 0, nil
}