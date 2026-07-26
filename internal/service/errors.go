package service

import "errors"

var (
	//Authorization
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrCreateUser         = errors.New("failed to create user")

	//User
	ErrUserNotFound = errors.New("user not found")

	//Conversation
	ErrCannotCreateConversation = errors.New("cannot create conversation")
	ErrConversationWithYourself = errors.New("cannot create conversation with yourself")
	ErrConversationNotFound     = errors.New("conversation not found")

	//Message
	ErrMessageTextEmpty   = errors.New("message text cannot be empty")
	ErrMessageTextTooLong = errors.New("message text is too long")
)
