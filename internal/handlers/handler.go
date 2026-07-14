package handlers

import (
	"log/slog"

	"github.com/amat-andre/corporate_messenger/internal/service"
)

type Handler struct {
	services *service.Services
	logger *slog.Logger
}

// Структура ответа с токеном для хендлеров.
type TokenResponse struct {
	Token string `json:"token"`
}

func New(services *service.Services, logger *slog.Logger) *Handler {
	return &Handler{
		services: services,
		logger: logger,
	}
}