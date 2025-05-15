package logger

import (
	"context"
	"log/slog"
)

// DiscardHandler used for testing
type DiscardHandler struct{}

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}

func (h *DiscardHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

func (h *DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *DiscardHandler) WithGroup(name string) slog.Handler {
	return h
}
