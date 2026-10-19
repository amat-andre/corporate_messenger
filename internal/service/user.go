package service

import (
	"errors"
	"log/slog"

	"github.com/amat-andre/corporate_messenger/internal/db"
	"github.com/amat-andre/corporate_messenger/internal/dto"
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

func (s *UserService) CurrentUser(userID int64) (*dto.CurrentUser, error) {
	// пока без особой логики, просто возвращаем пользователя по ID
	user, err := s.db.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &dto.CurrentUser{
		ID:    user.ID,
		Login: user.Login,
	}, nil
}

func (s *UserService) SearchUsers(search string) (*dto.UsersResponse, error) {
	users, err := s.db.SearchUsers(search)
	if err != nil {
		return nil, err
	}

	resp := &dto.UsersResponse{
		Users: make([]dto.UserShortResponse, 0, len(users)),
	}

	for _, user := range users {
		resp.Users = append(resp.Users, dto.UserShortResponse{
			ID:    user.ID,
			Login: user.Login,
		})
	}

	return resp, nil
}
