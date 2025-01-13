package palikka

import (
	"context"
	"log/slog"
)

type contextKey int

const (
	loggerContextKey  contextKey = iota
	sessionContextKey contextKey = iota
	userContextKey    contextKey = iota
)

func CopyContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func CopyContextWithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, session)
}

func CopyContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, _ := ctx.Value(loggerContextKey).(*slog.Logger)
	return logger
}

func SessionFromContext(ctx context.Context) *Session {
	session, _ := ctx.Value(sessionContextKey).(*Session)
	return session
}

func UserFromContext(ctx context.Context) *User {
	user, _ := ctx.Value(userContextKey).(*User)
	return user
}
