package auth

import (
	"fmt"
	"errors"

	"github.com/amat-andre/corporate_messenger/internal/models"
)

var ErrInvalidPassword = errors.New("invalid password")


// CheckCredentials сравнивает пароль вводимый пользователем с хэшем в базе.
func CheckCredentials(user *models.User, password string) error {
	if err := CheckPassword(user.PasswordHash, password); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidPassword, err)
	}

	return nil
}