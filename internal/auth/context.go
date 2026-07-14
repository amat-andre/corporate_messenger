package auth

import "context"

type contextKey string

const userIDKey contextKey = "userID"

// ContextWithUserID сохраняет ID пользователя в Context.
func ContextWithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext возвращает ID пользователя из Context.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}