package handlers

import (
	"log/slog"

	"github.com/amat-andre/corporate_messenger/internal/service"
)

type Handler struct {
	services *service.Services
	logger   *slog.Logger
}

func New(services *service.Services, logger *slog.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}
