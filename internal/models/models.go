package models

import "time"

// User представляет пользователя системы.
type User struct {
	ID           int64     `db:"id" json:"id"`
	Login        string    `db:"login" json:"login"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

func NewUser(login string, passwordHash string) *User {
	return &User{
		Login:        login,
		PasswordHash: passwordHash, // пока оставляем, посмотреть по логике
		CreatedAt:    time.Now(),
	}
}
