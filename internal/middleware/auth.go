package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/amat-andre/corporate_messenger/internal/auth"
	help "github.com/amat-andre/corporate_messenger/internal/helpers"
)

var (
	ErrMissingAuthorizationHeader = errors.New("authorization header is missing")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	// ErrInvalidToken               = errors.New("invalid token") - на будущее
	// и может ошибки авторизации можно вынести в отдельный файл errors.go пакета auth, чтобы не загромождать миддлварь
)

func Auth(next http.Handler, logger *slog.Logger,) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		
		header := strings.TrimSpace(req.Header.Get("Authorization"))

		if header == "" {
			help.WriteError(logger, w, http.StatusUnauthorized, ErrMissingAuthorizationHeader)
			return
		}

		const prefix = "Bearer "

		if !strings.HasPrefix(header, prefix) {
			help.WriteError(logger, w, http.StatusUnauthorized, ErrInvalidAuthorizationHeader)
			return
		}

		token := strings.TrimPrefix(header, prefix)

		userID, err := auth.ParseToken(token)
		if err != nil {
			help.WriteError(logger, w, http.StatusUnauthorized, err)
			return
		}

		ctx := auth.ContextWithUserID(req.Context(), userID)

		next.ServeHTTP(w, req.WithContext(ctx))
	})
}