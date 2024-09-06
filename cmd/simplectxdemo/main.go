package main

import (
	"context"
	"log/slog"

	"github.com/croepha/go-logging-extras/logctx"
)

func main() {

	// Setup base context, please also see also loginit package
	ctx := logctx.Context(context.Background(), slog.Default().Handler())

	// Optionally add an attribute to the context
	ctx = logctx.Attr(ctx, "GlobalAttr", "some_global_value")

	doSomething(ctx)
}

func doSomething(ctx context.Context) {

	// This logs:
	// 2024/09/06 12:52:23 INFO doing some stuff GlobalAttr=some_global_value Func=doSomething
	logctx.Info(ctx, "doing some stuff",
		"Func", "doSomething", // Optional attrs, just like slog.Logger.Info
		// ... more attrs
	)
}
