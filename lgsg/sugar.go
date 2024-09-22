package lgsg

import (
	"context"
	"log/slog"
	"slices"
	"sync"

	"github.com/croepha/go-logging-extras/errordump"
	"github.com/croepha/go-logging-extras/logctx"
	"github.com/croepha/go-logging-extras/loginit"
	"github.com/croepha/go-logging-extras/logwrap"
)

/*

This package provides some alternative syntax sugar for using slog

*/

var autoInitOnce sync.Once

// Function that returns a sugar handle and ensures that
// default init was done
func New() Sugar {

	autoInitOnce.Do(func() {
		logctx.Handler(loginit.MustInit(context.Background()))
	})

	return Sugar{}
}

// Sugar handle
type L struct {
	parent    *L
	attr      slog.Attr
	wrapDepth int
}

// TODO: as an alternative to single source, kinda would be
// nice to have a method to add full stack... probably
// should be done as a valuer

// Increments the wrap depth, which causes more frames to be skipped when resolving the source lines (only works when directly logging)
func (l L) Wrap() L {
	l.wrapDepth++
	return l
}

// Disables source code location resolution (only works when directly logging)
func (l L) NoSource() L {
	l.wrapDepth = logwrap.WrapDepth__DisablePC
	return l
}

func (l *L) resolveAttrs() []slog.Attr {
	attrs := make([]slog.Attr, 0, 10)
	for l != nil {
		attrs = append(attrs, l.attr)
		l = l.parent
	}
	slices.Reverse(attrs)
	return attrs
}

// Add attribute
func (l L) A(name string, value any) L {
	return L{parent: &l, attr: slog.Any(name, value), wrapDepth: l.wrapDepth}
}

// Adds an error attribute to be deeply described
func (l L) Err(err error) L {
	return L{parent: &l, attr: errordump.NewSlog("error", err), wrapDepth: l.wrapDepth}
}

// Addes attribute
func (l L) Attr(attr slog.Attr) L {
	return L{parent: &l, attr: attr, wrapDepth: l.wrapDepth}
}

// Gives new handler with given attributes
func (l L) Handler(ctx context.Context) slog.Handler {
	return logctx.Handler(ctx).WithAttrs(l.resolveAttrs())
}

// Gives a new ctx that has a handler with these attributes
func (l L) Context(ctx context.Context) context.Context {
	return logctx.Context(ctx, l.Handler(ctx))
}

// handles a new log record
func (l L) LogFromWrapper(ctx context.Context, additionalWrapDepth int, level slog.Level, msg string) {
	logwrap.LogAttrs(ctx, logctx.Handler(ctx), l.wrapDepth+additionalWrapDepth+1, level, l.resolveAttrs(), msg)
}

// handles a new log record
func (l L) Log(ctx context.Context, level slog.Level, msg string) {
	l.LogFromWrapper(ctx, 1, level, msg)
}

// handles a new Debug log record
func (l L) Debug(ctx context.Context, msg string) {
	l.LogFromWrapper(ctx, 1, slog.LevelDebug, msg)
}

// handles a new Info log record
func (l L) Info(ctx context.Context, msg string) {
	l.LogFromWrapper(ctx, 1, slog.LevelInfo, msg)
}

// handles a new Warn log record
func (l L) Warn(ctx context.Context, msg string) {
	l.LogFromWrapper(ctx, 1, slog.LevelWarn, msg)
}

// handles a new Error log record
func (l L) Error(ctx context.Context, msg string) {
	l.LogFromWrapper(ctx, 1, slog.LevelError, msg)
}

type common struct {
	parent *common // Creates a chain of handles that store the attributes
	// Some basic benchmarks actually show that this linked list approach is faster than using a slice
	// The slice approach is perhaps a bit more complicated then you would expect, as we would need to clone
	// the slice each time a method is called in order to use the chaining syntax...

	attr      slog.Attr
	wrapDepth int
}

// Gives new handler with given attributes
func (sgh common) Handler(h slog.Handler) slog.Handler {
	return h.WithAttrs(sgh.resolveAttrs())
}

// Increments the wrap depth, which causes more frames to be skipped when resolving the source lines (only works when directly logging)
func (sgh common) Wrap() common {
	sgh.wrapDepth++
	return sgh
}

// Disables source code location resolution (only works when directly logging)
func (sgh common) NoSource() common {
	sgh.wrapDepth = logwrap.WrapDepth__DisablePC
	return sgh
}

func (sgh *common) resolveAttrs() []slog.Attr {
	attrs := make([]slog.Attr, 0, 10)
	for sgh != nil {
		attrs = append(attrs, sgh.attr)
		sgh = sgh.parent
	}
	slices.Reverse(attrs)
	return attrs
}

// Add attribute from any value
func (sgh common) A(name string, value any) common {
	return sgh.Attr(slog.Any(name, value))
}

// Adds an error attribute to be deeply described
func (sgh common) errMaybe(err error) common {
	if err != nil {
		return sgh.Attr(errordump.NewSlog("error", err))
	}
	return sgh
}

// Addes slog.Attr attributes
func (sgh common) Attr(attrs ...slog.Attr) common {
	for _, attr := range attrs {
		bhClone := sgh // Force new copy to have new address
		sgh.parent = &bhClone
		sgh.attr = attr
	}
	return sgh
}

// handles a new log record, this provides all the functionality, able to be used in wrapping log calls
func (sh common) LogFromWrapper(ctx context.Context, sgh slog.Handler, additionalWrapDepth int, level slog.Level, msg string) {
	logwrap.LogAttrs(ctx, sgh, sh.wrapDepth+additionalWrapDepth+1, level, sh.resolveAttrs(), msg)
}

// Like Sugar, but it retains a handler instance instead and its logging methods do not take a Context
type Bound struct {
	c  common
	sh slog.Handler
}

// Base sugar handle
type Sugar struct {
	c common
}

// Gives new handler with given attributes
func (sgh Bound) Handler() slog.Handler { return sgh.c.Handler(sgh.sh) }

// Gives new handler with given attributes
func (sgh Sugar) Handler(ctx context.Context) slog.Handler { return sgh.c.Handler(logctx.Handler(ctx)) }

// Gives a new ctx that has a handler with these attributes
func (sgh Sugar) Context(ctx context.Context) context.Context {
	return logctx.Context(ctx, sgh.Handler(ctx))
}

// Gives a new ctx that has a handler with these attributes
func (sgh Bound) Context(ctx context.Context) context.Context {
	return logctx.Context(ctx, sgh.sh)
}

// Returns a sugar handle that has a handler from ctx bound to it
func (sgh Sugar) Bound(ctx context.Context) Bound { return Bound{c: sgh.c, sh: sgh.Handler(ctx)} }

// handles a new log record, this provides all the functionality, able to be used in wrapping log calls
func (sgh Bound) LogFromWrapper(additionalWrapDepth int, level slog.Level, msg string) {
	sgh.c.LogFromWrapper(context.Background(), sgh.sh, 1, level, msg)
}

// handles a new log record, this provides all the functionality, able to be used in wrapping log calls
func (sgh Sugar) LogFromWrapper(ctx context.Context, additionalWrapDepth int, level slog.Level, msg string) {
	sgh.c.LogFromWrapper(ctx, logctx.Handler(ctx), 1, level, msg)
}

func (sgh Bound) com(c common) Bound       { sgh.c = c; return sgh }
func (sgh Sugar) com(c common) Sugar       { sgh.c = c; return sgh }
func (sgh Bound) errMaybe(err error) Bound { return sgh.com(sgh.c.errMaybe(err)) }
func (sgh Sugar) errMaybe(err error) Sugar { return sgh.com(sgh.c.errMaybe(err)) }

// Increments the wrap depth, which causes more frames to be skipped when resolving the source lines (only works when directly logging)
func (sgh Bound) Wrap() Bound { return sgh.com(sgh.c.Wrap()) }

// Increments the wrap depth, which causes more frames to be skipped when resolving the source lines (only works when directly logging)
func (sgh Sugar) Wrap() Sugar { return sgh.com(sgh.c.Wrap()) }

// Disables source code location resolution (only works when directly logging)
func (sgh Bound) NoSource() Bound { return sgh.com(sgh.c.NoSource()) }

// Disables source code location resolution (only works when directly logging)
func (sgh Sugar) NoSource() Sugar { return sgh.com(sgh.c.NoSource()) }

// Add attribute from any value
func (sgh Bound) A(name string, value any) Bound { return sgh.com(sgh.c.A(name, value)) }

// Add attribute from any value
func (sgh Sugar) A(name string, value any) Sugar { return sgh.com(sgh.c.A(name, value)) }

// Adds slog.Attr attributes
func (sgh Bound) Attr(attrs ...slog.Attr) Bound { return sgh.com(sgh.c.Attr(attrs...)) }

// Adds slog.Attr attributes
func (sgh Sugar) Attr(attrs ...slog.Attr) Sugar { return sgh.com(sgh.c.Attr(attrs...)) }

// handles a new log record
func (sgh Bound) Log(level slog.Level, msg string) { sgh.LogFromWrapper(1, level, msg) }

// handles a new log record
func (sgh Sugar) Log(ctx context.Context, level slog.Level, msg string) {
	sgh.LogFromWrapper(ctx, 1, level, msg)
}

// handles a new Debug log record
func (sgh Bound) Debug(msg string) { sgh.LogFromWrapper(1, slog.LevelDebug, msg) }

// handles a new Debug log record
func (sgh Sugar) Debug(ctx context.Context, msg string) {
	sgh.LogFromWrapper(ctx, 1, slog.LevelDebug, msg)
}

// handles a new Info log record
func (sgh Bound) Info(msg string) { sgh.LogFromWrapper(1, slog.LevelInfo, msg) }

// handles a new Info log record
func (sgh Sugar) Info(ctx context.Context, msg string) {
	sgh.LogFromWrapper(ctx, 1, slog.LevelInfo, msg)
}

// handles a new Warn log record
func (sgh Bound) Warn(msg string) { sgh.LogFromWrapper(1, slog.LevelWarn, msg) }

// handles a new Warn log record
func (sgh Sugar) Warn(ctx context.Context, msg string) {
	sgh.LogFromWrapper(ctx, 1, slog.LevelWarn, msg)
}

// handles a new Error log record, if err is not nil, adds it as an error
func (sgh Bound) Error(err error, msg string) {
	sgh.errMaybe(err).LogFromWrapper(1, slog.LevelError, msg)
}

// handles a new Error log record, if err is not nil, adds it as an error
func (sgh Sugar) Error(ctx context.Context, err error, msg string) {
	sgh.errMaybe(err).LogFromWrapper(ctx, 1, slog.LevelError, msg)
}

// logs and panics if err is not nil
func (sgh Bound) MustNotError(err error) {
	if err != nil {
		sgh.errMaybe(err).LogFromWrapper(1, slog.LevelError, "panicing due to unexpected error")
	}
}

// logs and panics if err is not nil
func (sgh Sugar) MustNotError(ctx context.Context, err error) {
	if err != nil {
		sgh.errMaybe(err).LogFromWrapper(ctx, 1, slog.LevelError, "panicing due to unexpected error")
	}
}
