package handlers

import (
	"net/http"

	help "github.com/amat-andre/corporate_messenger/internal/helpers"
)

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	help.WriteJSON(h.logger, w, http.StatusOK, map[string]string{
		"service": "Corporate Messenger API",
		"status":  "running",
	})
}