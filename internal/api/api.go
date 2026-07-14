package api

import (
	"net/http"
	"log/slog"

	hdl "github.com/amat-andre/corporate_messenger/internal/handlers"
	"github.com/amat-andre/corporate_messenger/internal/middleware"
)

func New(handler *hdl.Handler, logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	// Авторизация и регистрация
	mux.HandleFunc("POST /api/login", handler.Login)
	mux.HandleFunc("POST /api/register", handler.Register)

	// Получение информации о текущем пользователе
	mux.Handle("GET /api/me", middleware.Auth(http.HandlerFunc(handler.Me), logger))
	mux.Handle("GET /api/users", middleware.Auth(http.HandlerFunc(handler.SearchUsers), logger))

	// Просто проверка работы сервера - убрать в будущем
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Messenger API"))
		})

	return mux
}
