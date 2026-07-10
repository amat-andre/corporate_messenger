package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amat-andre/corporate_messenger/internal/config"
	"github.com/amat-andre/corporate_messenger/internal/db"
	"github.com/amat-andre/corporate_messenger/internal/server"
)

func main(){
	logger := slog.Default()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config:", "error", err.Error())
		os.Exit(1)
	}

	base := db.New()

	err = base.Connect(cfg.DB)
	if err != nil {
		logger.Error("failed to connect to database:", "error", err.Error())
		os.Exit(1)
	}
	defer base.Close()

	srv := server.New(cfg.Server.Port, logger)
	
	go func() {
		if err := srv.Run(); err != nil {
			logger.Error("failed to run server:", "error", err.Error())
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

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to shutdown server:", "error", err.Error())
	}

	logger.Info("application stopped")
}