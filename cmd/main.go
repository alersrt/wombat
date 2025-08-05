package main

import (
	"context"
	"flag"
	"os"
	"wombat/internal"
	"wombat/pkg"
)

func main() {
	mainCtx, mainCancelCauseFunc := context.WithCancelCause(context.Background())
	defer mainCancelCauseFunc(nil)

	conf := new(internal.Config)
	args := parseArgs(os.Args)

	err := conf.Init(args)
	pkg.Try(err)

	dbStorage := internal.NewDbStorage(conf.PostgreSQL.Url)
	forwardChannel := make(chan *internal.Message)
	telegramSource, err := internal.NewTelegramSource(conf.Telegram.Token, forwardChannel, dbStorage)
	pkg.Try(err)
	jiraTarget := internal.NewJiraTarget(conf.Jira.Url, conf.Bot.Tag, dbStorage, forwardChannel)

	dmn := pkg.Create(conf)
	go func() {
		err := dmn.
			AddTask(jiraTarget.Do).
			AddTask(telegramSource.Do).
			Start(mainCtx)
		pkg.Try(err)
	}()
	select {}
}

func parseArgs(args []string) []string {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("config", "./cmd/config.yaml", "path to config")

	if err := flags.Parse(args[1:]); err != nil {
		panic(err)
	}

	return []string{*configPath}
}
