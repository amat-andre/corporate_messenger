package handlers

import (
	"net/http"
	"strings"

	help "github.com/amat-andre/corporate_messenger/internal/helpers"
)

func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {

	search := strings.TrimSpace(
		r.URL.Query().Get("search"),
	)

	users, err := h.services.Auth.SearchUsers(search)
	if err != nil {
		help.WriteError(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	help.WriteJSON(h.logger, w, http.StatusOK, users)
}