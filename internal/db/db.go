package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/amat-andre/corporate_messenger/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	DB *sql.DB
}

func New() *Database {
	return &Database{}
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
		return err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	base.DB = db

	return nil
}

func (base *Database) Close() error {
	if base.DB == nil {
		return nil
	}
	return base.DB.Close()
}