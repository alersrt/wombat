package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"wombat/internal"
	"wombat/pkg/daemon"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	app := new(internal.App)
	init := func() error {
		args, err := parseArgs()
		if err != nil {
			return err
		}
		err = app.Init(ctx, args)
		if err != nil {
			return err
		}
		return nil
	}
	shutdown := func() {
		app.Shutdown(cancel)
	}

	if err := init(); err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

	go app.Do(ctx)

	code, err := daemon.HandleSignals(ctx, shutdown, init)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
	os.Exit(code)
}

// parseArgs parses os.Args and returns list of parsed values. Here is descriptions by indexes:
// 0 - path to config file.
func parseArgs() ([]string, error) {
	args := os.Args
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("config", "./cmd/config.yaml", "path to config")

	err := flags.Parse(args[1:])
	if err != nil {
		return nil, fmt.Errorf("parseArgs: %w", err)
	}

	return []string{*configPath}, nil
}
