package logtest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
 General logging test helpers
*/

// This returns a struct that helps unit test logging wrappers
// the member H is a Handler that can be used as a slog.Handler
// See example usage in TestTestHandler
func NewTestHandler(t *testing.T) *testHandler {
	th := &testHandler{t: t}
	th.H = slog.NewJSONHandler(
		&th.buf, &slog.HandlerOptions{
			AddSource: true,
		},
	)
	th.reset()
	return th
}

func (th *testHandler) reset() {
	th.buf.Reset()
	th.dec = json.NewDecoder(&th.buf)
}

func (th *testHandler) RequireEOF() {
	tok, err := th.dec.Token()
	require.Equal(th.t, nil, tok)
	require.Equal(th.t, io.EOF, err)
	th.reset()
}

type testHandler struct {
	H   slog.Handler
	buf bytes.Buffer
	t   *testing.T
	dec *json.Decoder
}

// Requires that a specific line was logged
// Source line is expected to be previous line
// Also requires that this is the final line
func (th *testHandler) RequireLine(expectedLevel slog.Level, expectedMsg string, expectedArgs ...any) {
	th.t.Helper()
	th.RequireLineExtra(-1, 1, expectedLevel, expectedMsg, expectedArgs...)
	th.RequireEOF()
}

// Similar to RequireLine, but with more options and allows for multiple lines
// if sourceSkip is negative then:
//   if sourceLineDelta == 0 then source is expected to be omitted from output
//   if sourceLineDelta == 1 then source is ignored

// else sourceSkip says the number of frames to skip to get expected source line info
// sourceLineDelta is the number of lines to adjust for expected source line.  -1 is the preceeding line
func (th *testHandler) RequireLineExtra(sourceLineDelta int, sourceSkip int, expectedLevel slog.Level, expectedMsg string, expectedArgs ...any) {
	th.t.Helper()

	type jsonObject = map[string]any

	expectedBuf := bytes.Buffer{}
	slog.New(slog.NewJSONHandler(&expectedBuf, nil)).
		Log(context.Background(), expectedLevel, expectedMsg, expectedArgs...)
	expected := jsonObject{}
	require.NoError(th.t, json.Unmarshal(expectedBuf.Bytes(), &expected))
	delete(expected, "time")

	if sourceSkip >= 0 {
		_, sourcePath, sourceLine, ok := runtime.Caller(sourceSkip + 1)
		require.True(th.t, ok)
		expected["source"] = jsonObject{
			"file": sourcePath,
			"line": float64(sourceLine + sourceLineDelta),
		}
	}

	actual := jsonObject{}
	require.NoError(th.t, th.dec.Decode(&actual))

	delete(actual, "time")

	if sourceSkip < 0 && sourceLineDelta == 1 {
		delete(actual, "source")
	} else {
		if s := actual["source"]; s != nil {
			delete(s.(jsonObject), "function")
		}
	}

	require.Equal(th.t, expected, actual)
}
