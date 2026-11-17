package handlers

import (
	"errors"
	"net/http"
	"strings"

	help "github.com/amat-andre/corporate_messenger/internal/helpers"
)

var ErrEmptySearch = errors.New("search query is required")

func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {

	search := strings.TrimSpace(
		r.URL.Query().Get("search"),
	)

	if search == "" {
		help.WriteError(h.logger, w, http.StatusBadRequest, ErrEmptySearch)
		return
	}

	users, err := h.services.User.SearchUsers(search)
	if err != nil {
		help.WriteError(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	help.WriteJSON(h.logger, w, http.StatusOK, users)
}
