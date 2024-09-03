package logctx

import (
	"context"
	"log/slog"

	"github.com/croepha/go-logging-extras/logwrap"
)

/*

Some tools to make it easy to ctx based logging
A Slog handler can be added to and retrieved from the givent context

Debug/Info/Warn/Error are provided as convenience functions that use
Handler from ctx

*/

// NOTE: Intentially not using a custom type to avoid collisions
// We are using a normal string, with a known but likely to be unique value
// so that we can intentially allow other packages to use this value.
// So in-theory, we could have multiple versions or copies of this module imported by
// different modules, and still be compatible with eachother
var contextKey = "slog.Handler-7263656f68700a61"

// Gets Handler from context (or use the default if one isn't set)
func Handler(ctx context.Context) slog.Handler {
	v, _ := ctx.Value(contextKey).(slog.Handler)
	if v == nil {
		v = slog.Default().Handler()
	}
	return v
}

// Creates a new context with the given handler added to it
func Context(ctx context.Context, handler slog.Handler) context.Context {
	//lint:ignore SA1029 see contextKey comment
	return context.WithValue(ctx, contextKey, handler)
}

// Returns context with added slog attribute
func Attr(ctx context.Context, name string, value any) context.Context {
	return Context(ctx,
		Handler(ctx).WithAttrs(
			[]slog.Attr{
				slog.Any(name, value),
			},
		),
	)
}

// Log a Debug record using handler from context
func Debug(ctx context.Context, msg string, args ...any) {
	logwrap.Log(ctx, Handler(ctx), 1, slog.LevelDebug, args, msg)
}

// Log an Info record using handler from context
func Info(ctx context.Context, msg string, args ...any) {
	logwrap.Log(ctx, Handler(ctx), 1, slog.LevelInfo, args, msg)
}

// Log a Warn record using handler from context
func Warn(ctx context.Context, msg string, args ...any) {
	logwrap.Log(ctx, Handler(ctx), 1, slog.LevelWarn, args, msg)
}

// Log an Error record using handler from context
func Error(ctx context.Context, msg string, args ...any) {
	logwrap.Log(ctx, Handler(ctx), 1, slog.LevelError, args, msg)
}
