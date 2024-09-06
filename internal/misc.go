package internal

import (
	"context"
	"log/slog"
)

func (h *NullHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *NullHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *NullHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *NullHandler) WithGroup(name string) slog.Handler {
	return h
}

type NullHandler struct {
}
