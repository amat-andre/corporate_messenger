package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/amat-andre/corporate_messenger/internal/auth"
	"github.com/amat-andre/corporate_messenger/internal/dto"
	help "github.com/amat-andre/corporate_messenger/internal/helpers"
	"github.com/amat-andre/corporate_messenger/internal/service"
)

func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	conversationID, err := strconv.ParseInt(r.PathValue("conversationID"), 10, 64)
	if err != nil || conversationID <= 0 {
		help.WriteError(h.logger, w, http.StatusBadRequest, ErrInvalidConversationID)
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		help.WriteError(h.logger, w, http.StatusUnauthorized, auth.ErrUserUnauthorized)
		return
	}

	var req dto.CreateMessageRequest

	if err := help.ReadJSON(r, &req); err != nil {
		help.WriteError(h.logger, w, http.StatusBadRequest, err)
		return
	}

	message, err := h.services.Message.CreateMessage(conversationID, userID, req.Text)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMessageTextEmpty):
			help.WriteError(h.logger, w, http.StatusBadRequest, err)

		case errors.Is(err, service.ErrMessageTextTooLong):
			help.WriteError(h.logger, w, http.StatusBadRequest, err)

		case errors.Is(err, service.ErrConversationNotFound):
			help.WriteError(h.logger, w, http.StatusNotFound, err)

		default:
			help.WriteError(h.logger, w, http.StatusInternalServerError, err)
		}
		return
	}

	response := dto.MessageResponse{
		ID:             message.ID,
		UserID:         message.UserID,
		ConversationID: message.ConversationID,
		Username:       message.Username,
		Text:           message.Text,
		CreatedAt:      message.CreatedAt,
	}

	help.WriteJSON(h.logger, w, http.StatusCreated, response)
}

func (h *Handler) GetConversationMessages(w http.ResponseWriter, r *http.Request) {
	conversationID, err := strconv.ParseInt(r.PathValue("conversationID"), 10, 64)
	if err != nil || conversationID <= 0 {
		help.WriteError(h.logger, w, http.StatusBadRequest, ErrInvalidConversationID)
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		help.WriteError(h.logger, w, http.StatusUnauthorized, auth.ErrUserUnauthorized)
		return
	}

	messages, err := h.services.Message.GetConversationMessages(conversationID, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrConversationNotFound):
			help.WriteError(h.logger, w, http.StatusNotFound, err)

		default:
			help.WriteError(h.logger, w, http.StatusInternalServerError, err)
		}
		return
	}

	response := dto.MessagesResponse{
		Messages: make([]dto.MessageResponse, 0, len(messages)),
	}

	for _, message := range messages {
		response.Messages = append(response.Messages,
			dto.MessageResponse{
				ID:             message.ID,
				UserID:         message.UserID,
				ConversationID: message.ConversationID,
				Username:       message.Username,
				Text:           message.Text,
				CreatedAt:      message.CreatedAt,
			},
		)
	}

	help.WriteJSON(h.logger, w, http.StatusOK, response)
}
