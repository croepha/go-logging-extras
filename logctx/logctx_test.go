package logctx_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/croepha/go-logging-extras/logctx"
	"github.com/croepha/go-logging-extras/logtest"
)

func Test(t *testing.T) {
	ctx := context.Background()
	th := logtest.NewTestHandler(t)
	handler := th.H

	ctx = logctx.Context(ctx, handler)

	logctx.Info(ctx, "info test", "attr0", "foo")
	th.RequireLine(slog.LevelInfo, "info test", "attr0", "foo")
}
