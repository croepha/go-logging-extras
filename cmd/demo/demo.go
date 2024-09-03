package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/croepha/go-logging-extras/lgsg"
	"github.com/croepha/go-logging-extras/logctx"
	"github.com/croepha/go-logging-extras/loginit"
)

// This is just a shortcut, provides alternative convenient syntax sugar
var l lgsg.L

func main() {

	// Will setup some global state based on envs
	// example set SLOG_OUTPUT=/tmp/out.log to write to file instead of output
	// defaults to JSON logs if output isn't a terminal
	ctx := loginit.MustInit(context.Background())

	l.Debug(ctx, "only present with env SLOG_LEVEL=debug set")

	// Some global context, available everywhere
	ctx = l.A("WorkerID", "worker0").Context(ctx)

	l.Info(ctx, "Starting up main process") // includes WorkerID

	tasksToDo := []string{"task-10", "task-20", "task-30"}

	for idx, taskId := range tasksToDo {

		ctx := l.A("TaskID", taskId).Context(ctx)

		l.A("IDX", idx).Info(ctx, "adding attributes can be done through chaining") // includes WorkerID, TaskID and IDX

		l.Info(ctx, "Starting task") // includes WorkerID and TaskID

		doTask(ctx)

	}

	// Error dumping:
	if err := doSomeOperationThatFails(); err != nil {

		// includes WorkerID, also include a deep unwrapping of err that includes details
		// such the name of the package and type of the error, and any marshallable details
		// For example, when log output is set to a file, the path, op ("stat"), and errno (2)
		// are included as a structured attribute
		l.Err(err).Error(ctx, "doSomeOperationThatFails")
	}

}

// Alternatively, you could also preload a common attribute and it would keep that state everywhere it's used
// Combining retained attributes in the sugar handle, and also in the CTX gives two dimensions of attribute reuse
// which is powerful, for example some attribute follow the request or current task, and should be added to the context
// while other attributes follow a service or package
var l2 = l.A("Package", "demo")

func doTask(ctx context.Context) {

	l2.Info(ctx, "working on task") // includes WorkerID, TaskID and Package

	// If you don't like the sugar shorthand, there is also package level Methods that work
	logctx.Info(ctx, "This also works") // includes WorkerID, TaskID

	// Also, since we ran loginit.Init, the compatibility handler is installed so existing code will also work as long as it include the context:
	slog.InfoContext(ctx, "and this too") // includes WorkerID, TaskID

}

func doSomeOperationThatFails() error {
	_, err := os.Stat("/file/that/does/not/exist")
	return err

}
