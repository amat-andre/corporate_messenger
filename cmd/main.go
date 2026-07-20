package main

import (
	"context"
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
)

func main(){
	logger := slog.Default()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config:", "error", err)
		os.Exit(1)
	}
	auth.Init(cfg.JWT.Secret)

	database := db.New(logger)

	err = database.Connect(cfg.DB)
	if err != nil {
		logger.Error("failed to connect database:", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	err = database.Migrate()
	if err != nil {
		logger.Error("failed to migrate database:", "error", err)
		os.Exit(1)
	}

	services := service.NewServices(database, logger)

	handler := handlers.New(services, logger) // в будущем разбить на отдельные структуры, чтобы не передавать все сервисы в каждый хендлер

	router := api.New(handler, logger)

	server := server.New(cfg.Server.Port, router, logger)
	
	go func() {
		if err := server.Run(); err != nil {
			logger.Error("failed to run server:", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to shutdown server:", "error", err)
	}

	logger.Info("application stopped")
}