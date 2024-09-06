
# go-logging-extras

This module provides some enhancements to the standard log/slog package

## Context based logging

The `logctx` package allows you to use the `context.Context` value you already have plumbed throughout your code to retain request scoped logging attributes.

Example: (`cmd/simplectxdemo/main.go`)

    // Setup base context, please also see also loginit package
    ctx := logctx.Context(context.Background(), slog.Default().Handler())

    // Optionally add an attribute to the context
    ctx = logctx.Attr(ctx, "GlobalAttr", "some_global_value")

	// This logs:
	// 2024/09/06 12:52:23 INFO doing some stuff GlobalAttr=some_global_value Func=doSomething
	logctx.Info(ctx, "doing some stuff",
		"Func", "doSomething", // Optional attrs, just like slog.Logger.Info
		// ... more attrs
	)

This stores the given handler in the returned context.  The handler retains a set of attributes.

There are functions for each standard log level: `Debug`, `Info`, `Warn`, `Error`.  You can also get the Handler and use it directly with the `Handler` function.

You can also set `logctx.DefaultHandler` to a handler that will be used if there is not a handler set on a given context.  **IMPORTAINT**: You must set either the `DefaultHandler` or always have a `Handler` set or logs will be silently dropped.   Alternatively you can force a panic by setting `PanicOnNullHandler` to `true`.

There is also the `ctxhandler` package which provides a `slog.handler` implementation that can be used to allow existing code that is designed to just use a `slog.Logger` or `slog.Handler` directly to use the context based logging.

Example:

    slog.SetDefault(slog.New(ctxhandler.NewHandler()))

This code makes it so that any `slog.DebugContext(ctx, ...)` call will use the handler from the given context.

## Easy setup for logging

The `loginit` package provides a one simple function to do some sensible setup:
  - Configures a handler which is configurable via environment variables
  - Sets up and returns a context with the handler
  - Sets logctx.DefaultHandler
  - Installs the compatibility handler as default for slog

## Detailed error dumping

The `errordump` package provides some tools to inspect error objects and use them with structured logging.
Along with type specific details you can also easly determine specific type of a given error.

Example:

    func doSomeOperationThatFails() error {
        _, err := os.Stat("/file/that/does/not/exist")
        return err
    }

    func main() {
        ctx := loginit.MustInit(context.Background())
        err := doSomeOperationThatFails()
        logctx.Error(ctx, "doSomeOperationThatFails", errordump.NewSlog("error", err))
    }

Ran with json output enabled, gives the has the error attribute like this:

    "error": {
        "Error": {
            "String": "stat /file/that/does/not/exist: no such file or directory",
            "ReflectedName": "PathError",
            "ReflectedPackagePath": "io/fs",
            "NextDetails": {
                "Op": "stat",
                "Path": "/file/that/does/not/exist",
                "Err": 2
            }
        },
        "WrappedError": {
            "String": "no such file or directory",
            "ReflectedName": "Errno",
            "ReflectedPackagePath": "syscall",
            "NextDetails": 2
        }
    }

This is sometimes perfered over the default behavior:

    "error": "stat /file/that/does/not/exist: no such file or directory"

This has built-in functionality for unwrapping errors and reporting details deeply burried in a complex
error instance.  It is also configurable in a pluggable way by overwritting `GlobalDetailer`.  You can
implement your own `Detailer` to add details specific to the errors provided by the APIs you are using.
Detailers can be written in a way where they are composable and reusable.

## Installation

    go get github.com/croepha/go-logging-extras




