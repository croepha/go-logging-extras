package logwrap

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

/*

This package has some general helpers for logging

Specifically of interest is likely Log and LogAttrs which are alternatives to slog's Log/LogAttrs.
The Log and LogAttrs in this package are designed specifically to make writing log wrappers or
helper functions that log easier.  You can set the `wrapDepth` argument to control which frame is
used to get the source location information so that the filename and line number can be accurate.

*/

// Returns the program counter of the calling goroutine
// depth will skip a given number of stack frames
// in the case of func F() any { return PC(0) } will return a location inside of F
// program counter can be used with logging to translate to a source line
func PC(depth int) uintptr {
	var pcs [1]uintptr
	runtime.Callers(depth+2, pcs[:])
	return pcs[0]
}

const WrapDepth__DisablePC = -10000

// Create record
// wrapDepth will control how the PC (source line) is resolved
// wrapDepth defines the number of frames to skip
// if wrapDepth < 0 then PC is not resolved
// NOTE: Wrappers will add to wrapDepth without checking it's value, so
// you should set it to a large negative value, like WrapDepth__DisablePC to disable
func Record(wrapDepth int, level slog.Level, msg string) slog.Record {
	var pc uintptr
	if wrapDepth >= 0 {
		pc = PC(wrapDepth + 1)
	}

	return slog.NewRecord(
		time.Now(),
		level,
		msg,
		pc)

}

// Logs to given slog Handler
func Log(ctx context.Context, handler slog.Handler, wrapDepth int,
	level slog.Level, attrs []any, msg string) {
	if !handler.Enabled(ctx, level) {
		return
	}
	r := Record(wrapDepth+1, level, msg)
	r.Add(attrs...)
	err := handler.Handle(ctx, r)
	if err != nil {
		panic(err)
	}
}

// Logs to given slog Handler
func LogAttrs(ctx context.Context, handler slog.Handler, wrapDepth int,
	level slog.Level, attrs []slog.Attr, msg string) {
	if !handler.Enabled(ctx, level) {
		return
	}
	r := Record(wrapDepth+1, level, msg)
	r.AddAttrs(attrs...)
	err := handler.Handle(ctx, r)
	if err != nil {
		panic(err)
	}
}
