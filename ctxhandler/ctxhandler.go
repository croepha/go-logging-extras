package ctxhandler

import (
	"context"
	"log/slog"
	"slices"

	"github.com/croepha/go-logging-extras/logctx"
)

// Create a new handler instance
// This handler simply uses the supplied ctx to check for a Handler and uses it
// This handler can be used directly or installed as slog.Default()
// and provides compatibility for existing code that calls slog.InfoContext() or
// similar to use the logging handler configured in the ctx
func NewHandler() *ctxHandler {
	return &ctxHandler{}
}

type ctxHandler struct {
	attrs  []slog.Attr
	groups []string
}


func (ctxHandler) CannotBeLogCtxHandler() {}

func handler(ctx context.Context) slog.Handler {
	return logctx.Handler(ctx)
}

func (h *ctxHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return handler(ctx).Enabled(ctx, l)
}

func (h *ctxHandler) Handle(ctx context.Context, r slog.Record) error {
	real := handler(ctx)
	if len(h.attrs) > 0 {
		real = real.WithAttrs(h.attrs)
	}
	for _, g := range h.groups {
		real = real.WithGroup(g)
	}
	return real.Handle(ctx, r)
}

func (h *ctxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	r := *h
	r.attrs = slices.Concat(r.attrs, attrs)
	return &r
}

func (h *ctxHandler) WithGroup(name string) slog.Handler {
	r := *h
	r.groups = slices.Concat(r.groups, []string{name})
	return &r
}
