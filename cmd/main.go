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
	terminateIfError(err)

	err = conf.Init(args)
	terminateIfError(err)

	dbStorage, err := storage.NewDbStorage(conf.PostgreSQL.Url)
	terminateIfError(err)

	telegramSource, err := app.NewTelegramSource(conf.Telegram.Token)
	terminateIfError(err)

	dmn := daemon.Create(conf)

	runner, err := app.NewApplication(dmn, dbStorage, telegramSource)
	terminateIfError(err)

	runner.Run(mainCtx)
}

func terminateIfError(err error) {
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func parseArgs(args []string) ([]string, error) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("config", "./cmd/config.yaml", "path to config")

	if err := flags.Parse(args[1:]); err != nil {
		return nil, err
	}

	return []string{*configPath}, nil
}
