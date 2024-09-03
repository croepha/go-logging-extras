package logwrap_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/croepha/go-logging-extras/logtest"
	"github.com/croepha/go-logging-extras/logwrap"
)

func TestLog(t *testing.T) {
	ctx := context.Background()
	th := logtest.NewTestHandler(t)
	handler := th.H

	logwrap.Log(ctx, handler, 0, slog.LevelInfo, []any{"attr1", 10}, "test Log")
	th.RequireLine(slog.LevelInfo, "test Log", "attr1", 10)
}

func TestLogAttr(t *testing.T) {
	ctx := context.Background()
	th := logtest.NewTestHandler(t)
	handler := th.H

	logwrap.LogAttrs(ctx, handler, 0, slog.LevelInfo, []slog.Attr{slog.Any("attr1", 10)}, "test Log")
	th.RequireLine(slog.LevelInfo, "test Log", "attr1", 10)

}
