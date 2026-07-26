package handlers

import (
	"errors"
	"net/http"

	"github.com/amat-andre/corporate_messenger/internal/auth"
	"github.com/amat-andre/corporate_messenger/internal/dto"
	help "github.com/amat-andre/corporate_messenger/internal/helpers"
	"github.com/amat-andre/corporate_messenger/internal/service"
)

func (h *Handler) CreateConversation(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateDirectConversationRequest

	err := help.ReadJSON(r, &req)
	if err != nil {
		help.WriteError(h.logger, w, http.StatusBadRequest, err)
		return
	}

	currentUserID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		help.WriteError(h.logger, w, http.StatusUnauthorized, auth.ErrUserUnauthorized)
		return
	}

	conversation, err := h.services.Conversation.CreateDirectConversation(currentUserID, req.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrConversationWithYourself):
			help.WriteError(h.logger, w, http.StatusBadRequest, err)
		case errors.Is(err, service.ErrUserNotFound):
			help.WriteError(h.logger, w, http.StatusNotFound, err)
		default:
			help.WriteError(h.logger, w, http.StatusInternalServerError, err)
		}
		return
	}

	help.WriteJSON(h.logger, w, http.StatusCreated, conversation)
}
