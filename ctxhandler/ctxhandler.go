package ctxhandler

import (
	"context"
	"log/slog"
	"slices"

	"github.com/croepha/go-logging-extras/logctx"
)

var recurse slog.Handler

func init() {
	recurse = slog.Default().Handler()
}

// This sets the handler that is used by the NewHandler handler in the case where recursion is detected
// mutates global state without a lock, please serialize
func SetRecurseHandler(h slog.Handler) {
	recurse = h
}

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

func handler(ctx context.Context) slog.Handler {
	real := logctx.Handler(ctx)
	if _, ok := real.(*ctxHandler); ok {
		real = recurse
		slog.New(real).ErrorContext(ctx, "recursive use of handler abaited")
	}
	return real
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
