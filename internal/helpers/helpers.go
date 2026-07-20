package helpers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Структура ответа с ошибкой.
type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(logger *slog.Logger, w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

func ReadJSON(r *http.Request, target any) error {
	return json.NewDecoder(r.Body).Decode(target)
}

func WriteError(logger *slog.Logger, w http.ResponseWriter, statusCode int, err error) {
	WriteJSON(logger, w, statusCode, ErrorResponse{
		Error: err.Error(),
	})
} // зачем отдельный метод когда можно просто вызвать writeJSON с нужным статусом и ошибкой? может быть, чтобы не дублировать код, но тогда можно просто сделать функцию writeJSONWithError, которая будет принимать статус и ошибку и вызывать writeJSON с нужным статусом и ошибкой.
// потом сюда же засунуть readCookie() writeCookie() getBearerToken()