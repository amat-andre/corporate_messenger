package service

import (
	"log/slog"

	"github.com/amat-andre/corporate_messenger/internal/db"
)

type UserService struct {
	db     *db.Database
	logger *slog.Logger
}

func NewUserService(db *db.Database, logger *slog.Logger) *UserService {
	return &UserService{
		db:     db,
		logger: logger,
	}
}

type CurrentUser struct {
	ID    int64
	Login string
}

func (s *AuthService) CurrentUser(userID int64) (*CurrentUser, error) {
	// пока без особой логики, просто возвращаем пользователя по ID
	user, err := s.db.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &CurrentUser{
		ID:    user.ID,
		Login: user.Login,
	}, nil
}

type UserShortResponse struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

type UsersResponse struct {
	Users []UserShortResponse `json:"users"`
}

func (s *AuthService) SearchUsers(search string) (*UsersResponse, error) {
	users, err := s.db.SearchUsers(search)
	if err != nil {
		return nil, err
	}

	resp := &UsersResponse{
		Users: make([]UserShortResponse, 0, len(users)),
	}

	for _, user := range users {
		resp.Users = append(resp.Users, UserShortResponse{
			ID:    user.ID,
			Login: user.Login,
		})
	}

	return resp, nil
}