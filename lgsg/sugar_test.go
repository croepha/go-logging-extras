package lgsg_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/croepha/go-logging-extras/internal"
	"github.com/croepha/go-logging-extras/lgsg"
	"github.com/croepha/go-logging-extras/logctx"
	"github.com/croepha/go-logging-extras/logtest"
)

var l lgsg.L

func TestSugar(t *testing.T) {

	ctx := context.Background()

	th := logtest.NewTestHandler(t)

	handler := th.H

	ctx = logctx.Context(ctx, handler)
	l.A("attr1", 10).A("attr2", 20).Info(ctx, "message")
	th.RequireLine(slog.LevelInfo, "message", "attr1", 10, "attr2", 20)

	func() {
		l.Wrap().A("attr1", 10).A("attr2", 20).Info(ctx, "message")
	}()
	th.RequireLine(slog.LevelInfo, "message", "attr1", 10, "attr2", 20)

	l.NoSource().A("attr1", 10).A("attr2", 20).Info(ctx, "message")
	th.RequireLineExtra(0, -1, slog.LevelInfo, "message", "attr1", 10, "attr2", 20)
	th.RequireEOF()

}

func BenchmarkSugar(b *testing.B) {
	handler := &internal.NullHandler{}
	ctx := context.Background()
	ctx = logctx.Context(ctx, handler)
	l := lgsg.L{}
	for i := range b.N {
		l.A("bench_i", i).Info(ctx, "test line")
	}
}

func BenchmarkBaselineLogger(b *testing.B) {
	handler := &internal.NullHandler{}

	ctx := logctx.Context(context.Background(), handler)

	logger := slog.New(handler)
	for i := range b.N {
		logger.InfoContext(ctx, "test line", "bench_i", i)
	}
}

func BenchmarkBaselineAttrs(b *testing.B) {
	handler := &internal.NullHandler{}

	ctx := logctx.Context(context.Background(), handler)

	logger := slog.New(handler)
	for i := range b.N {
		logger.LogAttrs(ctx, slog.LevelInfo, "test line", slog.Int("bench_i", i))
	}
}
