package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var version = "dev"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	return runWithExecutor(args, executeScan)
}

func runWithExecutor(args []string, executor scanExecutor) int {
	root := newRootCommand(executor)
	ctx, stop := signal.NotifyContext(root.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	root.SetContext(ctx)
	root.SetArgs(args)

	exitCode := 0
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
	}
	if value, ok := root.Context().Value(exitCodeContextKey{}).(*int); ok && value != nil && *value != 0 {
		exitCode = *value
	}
	return exitCode
}
