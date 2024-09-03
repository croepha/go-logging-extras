package logtest_test

import (
	"log/slog"
	"testing"

	"github.com/croepha/go-logging-extras/logtest"
)

func TestTestHandler(t *testing.T) {

	th := logtest.NewTestHandler(t)

	slog.New(th.H).Info("test message", "attr0", "foo", "attr1", "bar")
	th.RequireLine(slog.LevelInfo, "test message", "attr0", "foo", "attr1", "bar")

}
