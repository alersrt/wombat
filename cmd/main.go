package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"wombat/internal"
)

const (
	ExitCodeDone         = 0
	ExitCodeError        = 1
	ExitCodeInvalidUsage = 2
	ExitCodeInterrupt    = 128 + int(syscall.SIGINT)
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ctx, cancelCause := context.WithCancelCause(ctx)

	go func() {
		cfg, err := internal.NewConfig(parseArgs())
		if err != nil {
			cancelCause(err)
		}

		kernel, err := internal.NewKernel(cfg)
		if err != nil {
			cancelCause(err)
		}

		kernel.Serve(ctx)
	}()

	<-ctx.Done()

	if ctx.Err() != nil {
		slog.Error(fmt.Sprintf("%+v", ctx.Err()))
		os.Exit(ExitCodeError)
	} else {
		slog.Info("finished")
		os.Exit(ExitCodeDone)
	}
}

// parseArgs parses os.Args and returns list of parsed values. Here is descriptions by indexes:
// 0 - path to config file.
func parseArgs() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config")
	flag.Parse()
	return path
}
