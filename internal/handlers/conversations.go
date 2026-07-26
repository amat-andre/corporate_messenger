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

var ErrInvalidConversationID = errors.New("invalid conversation ID")

func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		help.WriteError(h.logger, w, http.StatusUnauthorized, auth.ErrUserUnauthorized)
		return
	}

	conversations, err := h.services.Conversation.GetUserConversations(userID)
	if err != nil {
		help.WriteError(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	response := dto.ConversationsResponse{
		Conversations: make([]dto.ConversationResponse, 0, len(conversations)),
	}

	for _, conversation := range conversations {

		item := dto.ConversationResponse{
			ID:        conversation.ID,
			Type:      conversation.Type,
			CreatedAt: conversation.CreatedAt,
		}

		if conversation.OtherUser != nil {
			item.User = &dto.ConversationUserResponse{
				ID:    conversation.OtherUser.ID,
				Login: conversation.OtherUser.Login,
			}
		}

		response.Conversations = append(response.Conversations, item)
	}

	help.WriteJSON(h.logger, w, http.StatusOK, response)
}

func (h *Handler) GetConversationByID(w http.ResponseWriter, r *http.Request) {

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

	conversation, err := h.services.Conversation.GetConversationByID(conversationID, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrConversationNotFound):
			help.WriteError(h.logger, w, http.StatusNotFound, err)
		default:
			help.WriteError(h.logger, w, http.StatusInternalServerError, err)
		}
		return
	}

	response := dto.ConversationResponse{
		ID:        conversation.ID,
		Type:      conversation.Type,
		CreatedAt: conversation.CreatedAt,
	}

	if conversation.OtherUser != nil {
		response.User = &dto.ConversationUserResponse{
			ID:    conversation.OtherUser.ID,
			Login: conversation.OtherUser.Login,
		}
	}

	help.WriteJSON(h.logger, w, http.StatusOK, response)
}
