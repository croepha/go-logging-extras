package loginit

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/croepha/go-logging-extras/ctxhandler"
	"github.com/croepha/go-logging-extras/logctx"
	"golang.org/x/term"
)

// TODO: It would be nice if we could connect slog to t.Log when running tests...

// Perform some common startup things for slogs default logger
// mutates global state without a lock, please serialize
func Init(ctx context.Context) (context.Context, error) {

	handler, err := EnvHandler()
	if err != nil {
		return ctx, err
	}

	// Setup the ctx compatibility handler
	ctxhandler.SetRecurseHandler(handler)
	slog.SetDefault(slog.New(ctxhandler.NewHandler()))

	ctx = logctx.Context(ctx, handler)

	return ctx, nil
}

// Like Init, but panics on any errors
func MustInit(ctx context.Context) context.Context {
	ctx, err := Init(ctx)
	if err != nil {
		panic(err)
	}
	return ctx
}

// env SLOG_LEVEL configured default log level according to slog.Level.UnmarshalText
// examples: SLOG_LEVEL=debug SLOG_LEVEL=info
// env SLOG_OUTPUT sets the output
// it is set to a path that contains at-least one path separator or stdout or stderr
// other env vars starting with `SLOG_` may be used in the future
func EnvHandler() (slog.Handler, error) {
	var level slog.Level

	if e := os.Getenv("SLOG_LEVEL"); e != "" {
		if err := level.UnmarshalText([]byte(e)); err != nil {
			return nil, fmt.Errorf("SLOG_LEVEL: %+q unparsable: %w", e, err)
		}
	}

	var out io.Writer
	switch e := os.Getenv("SLOG_OUTPUT"); strings.ToLower(e) {
	// TODO syslog? or URLS and other protocols? Could have pluggable protocol handlers
	case "stderr", "":
		out = os.Stderr
	case "stdout":
		out = os.Stdout
	default:
		ps := string(os.PathSeparator)
		if !strings.Contains(e, ps) {
			return nil, fmt.Errorf("SLOG_OUTPUT: %+q should be a stderr,stdout or a path that contains %+q", e, ps)
		}
		f, err := os.OpenFile(e, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("SLOG_OUTPUT: %+q err: %w", e, err)
		}
		out = f
	}

	text := false
	if out, ok := out.(*os.File); ok && term.IsTerminal(int(out.Fd())) {
		text = true
	}

	opts := slog.HandlerOptions{
		Level: &level,
	}

	if text {
		return slog.NewTextHandler(out, &opts), nil
	}
	return slog.NewJSONHandler(out, &opts), nil

}
