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
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		app := new(internal.App)
		if err := app.Init(); err != nil {
			slog.Error(fmt.Sprintf("%+v", err))
		}

		slog.Info("start")
		err := app.Do(ctx)
		if err != nil {
			slog.Error(fmt.Sprintf("%+v", err))
		}
	}()

	<-ctx.Done()
	os.Exit(ExitCodeDone)
}

// parseArgs parses os.Args and returns list of parsed values. Here is descriptions by indexes:
// 0 - path to config file.
func parseArgs() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config")
	flag.Parse()
	return path
}
