package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"log/slog"
	"os"
	"wombat/internal"
	"wombat/pkg/daemon"
)

func main() {
	ctx := context.Background()

	app := new(internal.App)
	init := func() error {
		args, err := parseArgs(os.Args)
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
		os.Exit(1)
	}

	go app.Do(ctx)

	code, err := daemon.HandleSignals(ctx, app)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
	os.Exit(code)
}

func parseArgs(args []string) ([]string, error) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("config", "./cmd/config.yaml", "path to config")

	err := flags.Parse(args[1:])
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return []string{*configPath}, nil
}
