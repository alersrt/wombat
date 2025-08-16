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
		err = app.Init(args)
		if err != nil {
			return err
		}
		return nil
	}

	if err := init(); err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
		os.Exit(daemon.ExitCodeError)
	}

	go app.Do(ctx)

	code, err := daemon.HandleSignals(ctx, cancel)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
	slog.Info("exit:", "code", code)
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
