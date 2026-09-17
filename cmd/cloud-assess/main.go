package main

import (
	"context"
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
	root := newRootCommand(func(ctx context.Context, flags scanFlags) (int, error) {
		return executeScan(ctx, flags)
	})
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
