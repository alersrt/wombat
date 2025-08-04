package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"wombat/internal"
	"wombat/pkg/daemon"
)

func main() {
	mainCtx, mainCancelCauseFunc := context.WithCancelCause(context.Background())
	defer mainCancelCauseFunc(nil)

	conf := new(internal.Config)
	args, err := parseArgs(os.Args)
	terminateIfError(err)

	err = conf.Init(args)
	terminateIfError(err)

	dbStorage, err := internal.NewDbStorage(conf.PostgreSQL.Url)
	terminateIfError(err)

	forwardChannel := make(chan *internal.Message)

	telegramSource, err := internal.NewTelegramSource(conf.Telegram.Token, forwardChannel, dbStorage)
	terminateIfError(err)

	jiraTarget, err := internal.NewJiraTarget(conf.Jira.Url, conf.Bot.Tag, dbStorage, forwardChannel)
	terminateIfError(err)

	dmn := daemon.Create(conf)
	go dmn.
		AddTask(jiraTarget.Process).
		AddTask(telegramSource.Process).
		Start(mainCtx)

	select {}
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
