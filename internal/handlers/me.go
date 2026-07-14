package handlers

import (
	"errors"
	"net/http"

	"github.com/amat-andre/corporate_messenger/internal/auth"
	help "github.com/amat-andre/corporate_messenger/internal/helpers"
)

func (h *Handler) Me(w http.ResponseWriter, req *http.Request) {

	userID, ok := auth.UserIDFromContext(req.Context())

	if !ok {
		help.WriteError(h.logger, w, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	user, err := h.services.Auth.CurrentUser(userID)
	if err != nil {
		help.WriteError(h.logger, w, http.StatusInternalServerError, err) // cтатус 500, так как ошибка на сервере, а не в запросе, а ошибку потом поменять
		return
		}

	help.WriteJSON(h.logger, w, http.StatusOK, user)
}
