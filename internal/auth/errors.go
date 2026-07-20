package auth

import "errors"

// дописать ошибки для авторизации, либо удалить этот файл, либо удалить текущие ошибки (если они остаются в мидлваре)

var (
	ErrMissingAuthorizationHeader = errors.New("authorization header is missing")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	ErrInvalidToken               = errors.New("invalid token")
)