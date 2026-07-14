package handlers

import (
	"net/http"

	help "github.com/amat-andre/corporate_messenger/internal/helpers"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	help.WriteJSON(h.logger, w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}