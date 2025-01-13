package logs

import (
	"context"
	"github.com/stretchr/testify/mock"
	"log/slog"
)

type ContextualHandlerMock struct {
	mock.Mock
}

func (m *ContextualHandlerMock) Enabled(ctx context.Context, level slog.Level) bool {
	args := m.Called(ctx, level)
	return args.Bool(0)
}

func (m *ContextualHandlerMock) Handle(ctx context.Context, r slog.Record) error {
	args := m.Called(ctx, r)
	err, _ := args.Error(0).(error)
	return err
}

func (m *ContextualHandlerMock) WithAttrs(attrs []slog.Attr) slog.Handler {
	args := m.Called(attrs)
	return args.Get(0).(slog.Handler)
}

func (m *ContextualHandlerMock) WithGroup(name string) slog.Handler {
	args := m.Called(name)
	return args.Get(0).(slog.Handler)
}
