package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amat-andre/corporate_messenger/internal/api"
	"github.com/amat-andre/corporate_messenger/internal/auth"
	"github.com/amat-andre/corporate_messenger/internal/config"
	"github.com/amat-andre/corporate_messenger/internal/db"
	"github.com/amat-andre/corporate_messenger/internal/handlers"
	"github.com/amat-andre/corporate_messenger/internal/server"
	"github.com/amat-andre/corporate_messenger/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}

	slog.Info("application stopped")
}

func run() error {
	if err := godotenv.Load(); err != nil {
		slog.Debug(".env file not found, environment variables are used")
	}

	logger := slog.Default()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	auth.Init(cfg.JWT.Secret)

	database := db.New(logger)

	err = database.Connect(cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	defer database.Close()

	err = database.Migrate()
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	services := service.NewServices(database, logger)

	handler := handlers.New(services, logger) // в будущем разбить на отдельные структуры, чтобы не передавать все сервисы в каждый хендлер

	router := api.New(handler, logger)

	server := server.New(cfg.Server.Port, router, logger)

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- server.Run()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	select {
	// Ошибка запуска или неожиданное завершение сервера:
	// выходим с ошибкой, defer-ы (в т.ч. database.Close) отрабатывают.
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
		return nil

	// Получен сигнал остановки — плавно завершаем работу.
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
