package ctxhandler_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/croepha/go-logging-extras/ctxhandler"
	"github.com/croepha/go-logging-extras/logctx"
	"github.com/croepha/go-logging-extras/logtest"
)

func Test(t *testing.T) {
	ctx := context.Background()
	th := logtest.NewTestHandler(t)
	handler := th.H

	ctx = logctx.Context(ctx, handler)

	l := slog.New(ctxhandler.NewHandler())

	l.InfoContext(ctx, "info test", "attr0", "foo")
	th.RequireLine(slog.LevelInfo, "info test", "attr0", "foo")

	l.With("with0", "with0").InfoContext(ctx, "info test", "attr0", "foo")
	th.RequireLine(slog.LevelInfo, "info test", "with0", "with0", "attr0", "foo")

	l.WithGroup("withGroup0").InfoContext(ctx, "info test", "attr0", "foo")
	th.RequireLine(slog.LevelInfo, "info test", "withGroup0", map[string]any{"attr0": "foo"})

}

func TestRecursion(t *testing.T) {
	ctx := context.Background()
	th := logtest.NewTestHandler(t)

	ctxhandler.SetRecurseHandler(th.H)
	ctx = logctx.Context(ctx, ctxhandler.NewHandler())
	l := slog.New(ctxhandler.NewHandler())

	l.InfoContext(ctx, "info test", "attr0", "foo")
	th.RequireLineExtra(1, -1, slog.LevelError, "recursive use of handler abaited")
	th.RequireLineExtra(1, -1, slog.LevelError, "recursive use of handler abaited")
	th.RequireLineExtra(-3, 0, slog.LevelInfo, "info test", "attr0", "foo")
	th.RequireEOF()

}
