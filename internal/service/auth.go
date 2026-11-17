package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/amat-andre/corporate_messenger/internal/auth"
	"github.com/amat-andre/corporate_messenger/internal/db"
	"github.com/amat-andre/corporate_messenger/internal/models"
)

type AuthService struct {
	db     *db.Database
	logger *slog.Logger
}

func NewAuthService(db *db.Database, logger *slog.Logger) *AuthService {
	return &AuthService{
		db:     db,
		logger: logger,
	}
}

func (s *AuthService) Register(login, password string) (string, error) {
	// Создаем нового пользователя и хэшируем пароль.
	user, err := registerUser(login, password)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCreateUser, err)
	}

	// Сохраняем в БД. Уникальный индекс по login защищает
	// от гонки при параллельной регистрации.
	if err := s.db.CreateUser(user); err != nil {
		if db.IsUniqueViolation(err) {
			return "", ErrUserAlreadyExists
		}
		return "", fmt.Errorf("%w: %w", ErrCreateUser, err)
	}

	// Создаем JWT токен.
	token, err := auth.CreateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCreateUser, err)
	}

	s.logger.Info("user registered", "user_id", user.ID)

	return token, nil
}

func (s *AuthService) Login(login, password string) (string, error) {

	user, err := s.db.GetUserByLogin(login)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}

		return "", err
	}

	if err := auth.CheckCredentials(user, password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to login: %w", err)
	}

	s.logger.Info("user loggined", "user_id", user.ID)

	return token, nil
}

// registerUser создает модель нового пользователя и хэширует пароль.
func registerUser(login, password string) (*models.User, error) {
	// Хешируем пароль.
	passwordHash, err := auth.CreatePasswordHash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаем пользователя.
	user := models.NewUser(
		login,
		passwordHash,
	)

	return user, nil
}
