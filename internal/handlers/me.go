package handlers

import (
	"errors"
	"net/http"

	"github.com/amat-andre/corporate_messenger/internal/auth"
	help "github.com/amat-andre/corporate_messenger/internal/helpers"
	"github.com/amat-andre/corporate_messenger/internal/service"
)

func (h *Handler) Me(w http.ResponseWriter, req *http.Request) {

	userID, ok := auth.UserIDFromContext(req.Context())
	if !ok {
		help.WriteError(h.logger, w, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	user, err := h.services.User.CurrentUser(userID)
	if err != nil {
		// Пользователя нет в базе, хотя токен валиден — считаем запрос неавторизованным.
		if errors.Is(err, service.ErrUserNotFound) {
			help.WriteError(h.logger, w, http.StatusUnauthorized, err)
			return
		}
		help.WriteError(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	help.WriteJSON(h.logger, w, http.StatusOK, user)
}
