package logs

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

type ContextKey[T any] struct {
	LogKey string
}

func (k ContextKey[T]) Get(ctx context.Context) (T, error) {
	if v, ok := ctx.Value(k).(T); ok {
		return v, nil
	}
	var r T
	return r, errors.New("context key not found")
}

func (k ContextKey[T]) AddToContext(ctx context.Context, v T) context.Context {
	return context.WithValue(ctx, k, v)
}

func (k ContextKey[T]) AddToRecord(ctx context.Context, r slog.Record) slog.Record {
	if v, ok := ctx.Value(k).(T); ok {
		r.AddAttrs(slog.Any(k.LogKey, v))
	}
	return r
}

var (
	contextKeyRequestMethod = ContextKey[string]{LogKey: "method"}
	contextKeyRequestPath   = ContextKey[string]{LogKey: "path"}
	contextKeySessionID     = ContextKey[string]{LogKey: "session_id"}
)

func ContextKeyRequestMethod() ContextKey[string] {
	return contextKeyRequestMethod
}

func ContextKeyRequestPath() ContextKey[string] {
	return contextKeyRequestPath
}

func ContextKeySessionID() ContextKey[string] {
	return contextKeySessionID
}

func AddRequestAttributes(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = contextKeyRequestMethod.AddToContext(ctx, r.Method)
	ctx = contextKeyRequestPath.AddToContext(ctx, r.URL.Path)
	return r.WithContext(ctx)
}

type ContextualHandler interface {
	Enabled(ctx context.Context, level slog.Level) bool
	Handle(ctx context.Context, r slog.Record) error
	WithAttrs(attrs []slog.Attr) slog.Handler
	WithGroup(name string) slog.Handler
}

type contextualHandler struct {
	slog.Handler
}

func NewContextualHandler(h slog.Handler) ContextualHandler {
	return &contextualHandler{
		Handler: h,
	}
}

func (h *contextualHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *contextualHandler) Handle(ctx context.Context, r slog.Record) error {
	r = contextKeyRequestMethod.AddToRecord(ctx, r)
	r = contextKeyRequestPath.AddToRecord(ctx, r)
	r = contextKeySessionID.AddToRecord(ctx, r)

	// add other desired context keys here

	return h.Handler.Handle(ctx, r)
}

func (h *contextualHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.Handler.WithAttrs(attrs)
}

func (h *contextualHandler) WithGroup(name string) slog.Handler {
	return h.Handler.WithGroup(name)
}
