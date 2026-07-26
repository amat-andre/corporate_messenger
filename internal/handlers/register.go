package handlers

import (
	"errors"
	"net/http"

	"github.com/amat-andre/corporate_messenger/internal/dto"
	help "github.com/amat-andre/corporate_messenger/internal/helpers"
	"github.com/amat-andre/corporate_messenger/internal/service"
	"github.com/amat-andre/corporate_messenger/internal/validation"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := help.ReadJSON(r, &req); err != nil {
		help.WriteError(h.logger, w, http.StatusBadRequest, err)
		return
	}

	// проверка логина и пароля
	if err := validation.Credentials(req.Login, req.Password); err != nil {
		help.WriteError(h.logger, w, http.StatusBadRequest, err)
		return
	}

	// проверка пользователя в базе данных и получение токена при успешной регистрации
	token, err := h.services.Auth.Register(req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserAlreadyExists):
			help.WriteError(h.logger, w, http.StatusConflict, err)
		default:
			help.WriteError(h.logger, w, http.StatusInternalServerError, err)
		}
		return
	}

	help.WriteJSON(h.logger, w, http.StatusCreated, dto.TokenResponse{
		Token: token,
	})
}
