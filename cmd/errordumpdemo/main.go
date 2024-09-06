package main

import (
	"context"
	"os"

	"github.com/croepha/go-logging-extras/logctx"
	"github.com/croepha/go-logging-extras/loginit"
)

func doSomeOperationThatFails() error {
	_, err := os.Stat("/file/that/does/not/exist")
	return err
}

func main() {
	ctx := loginit.MustInit(context.Background())
	err := doSomeOperationThatFails()
	// logctx.Error(ctx, "doSomeOperationThatFails", errordump.NewSlog("error", err))
	logctx.Error(ctx, "doSomeOperationThatFails", "error", err)
}
