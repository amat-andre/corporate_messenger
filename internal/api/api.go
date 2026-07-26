package api

import (
	"log/slog"
	"net/http"

	hdl "github.com/amat-andre/corporate_messenger/internal/handlers"
	"github.com/amat-andre/corporate_messenger/internal/middleware"
)

func New(handler *hdl.Handler, logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	// Авторизация и регистрация
	mux.HandleFunc("POST /api/login", handler.Login)
	mux.HandleFunc("POST /api/register", handler.Register)

	// Получение информации о текущем пользователе
	mux.Handle("GET /api/me", middleware.Auth(handler.Me, logger))
	// Получение информации о существующих пользователях
	mux.Handle("GET /api/users", middleware.Auth(handler.SearchUsers, logger))

	// Создание личного диалога с указанным пользователем
	mux.Handle("POST /api/conversations", middleware.Auth(handler.CreateConversation, logger)) //над ручкой подумать
	// Получение списка диалогов текущего пользователя
	mux.Handle("GET /api/conversations", middleware.Auth(handler.GetConversations, logger))
	// Получение пользователем своего диалога по ID
	mux.Handle("GET /api/conversations/{conversationID}", middleware.Auth(handler.GetConversationByID, logger))
	// Отправка пользователем сообщения в указанный диалог
	mux.Handle("POST /api/conversations/{conversationID}/messages", middleware.Auth(handler.CreateMessage, logger))
	// Получение пользователем списка сообщений из указанного диалога
	mux.Handle("GET /api/conversations/{conversationID}/messages", middleware.Auth(handler.GetConversationMessages, logger))

	// Просто проверка работы сервера - убрать когда фронт будет
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Messenger API"))
	})

	return mux
}
