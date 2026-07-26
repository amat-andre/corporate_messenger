package validation

import (
	"errors"
	"strings"
)

var (
	ErrEmptyLogin    = errors.New("login is required")
	ErrEmptyPassword = errors.New("password is required")
	ErrShortPassword = errors.New("password must contain at least 8 characters")
	ErrLongLogin     = errors.New("login must not exceed 64 characters")
	ErrLongPassword  = errors.New("password must not exceed 128 characters")
)

func Credentials(login, password string) error {
	login = strings.TrimSpace(login)
	password = strings.TrimSpace(password)

	switch {
	case login == "":
		return ErrEmptyLogin

	case password == "":
		return ErrEmptyPassword

	case len(login) > 64:
		return ErrLongLogin

	case len(password) > 128:
		return ErrLongPassword

	case len(password) < 8:
		return ErrShortPassword
	}

	return nil
}
