package helpers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// maxBodySize ограничивает размер тела запроса (1 МБ).
const maxBodySize = 1 << 20

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

// ReadJSON декодирует тело запроса, ограничивая его размер.
func ReadJSON(w http.ResponseWriter, r *http.Request, target any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	return json.NewDecoder(r.Body).Decode(target)
}

// WriteError пишет ошибку в ответ. Для 5xx внутреннее сообщение
// попадает только в лог, клиент получает обобщенный текст.
func WriteError(logger *slog.Logger, w http.ResponseWriter, statusCode int, err error) {
	if statusCode >= http.StatusInternalServerError {
		logger.Error("internal error", "status", statusCode, "error", err)
		err = errors.New("internal server error")
	}

	WriteJSON(logger, w, statusCode, ErrorResponse{
		Error: err.Error(),
	})
}

// потом сюда же засунуть readCookie() writeCookie() getBearerToken()
