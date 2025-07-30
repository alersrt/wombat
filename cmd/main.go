package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"wombat/internal/app"
	"wombat/internal/storage"
	"wombat/pkg/daemon"
)

func main() {
	mainCtx, mainCancelCauseFunc := context.WithCancelCause(context.Background())
	defer mainCancelCauseFunc(nil)

	conf := new(app.Config)
	args, err := parseArgs(os.Args)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	err = conf.Init(args)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	aclRepository, err := storage.NewDbStorage(&conf.PostgreSQL.Url)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	telegram, err := app.NewTelegramSource(conf.Telegram.Token)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	dmn := daemon.Create(conf)

	runner, err := app.NewApplication(dmn, aclRepository, telegram)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	runner.Run(mainCtx)
}

func parseArgs(args []string) ([]string, error) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("config", "./cmd/config.yaml", "path to config")

	if err := flags.Parse(args[1:]); err != nil {
		return nil, err
	}

	return []string{*configPath}, nil
}
