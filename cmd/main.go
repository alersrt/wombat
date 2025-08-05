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
	pkg.Throw(err)

	db, err := internal.NewDbStorage(conf.PostgreSQL.Url)
	pkg.Throw(err)
	forwardChannel := make(chan *internal.Message)
	telegramSource, err := internal.NewTelegramSource(conf.Telegram.Token, forwardChannel, db)
	pkg.Throw(err)
	jiraTarget := internal.NewJiraTarget(conf.Jira.Url, conf.Bot.Tag, db, forwardChannel)

	dmn := pkg.Create(conf)
	go func() {
		err := dmn.
			AddTask(jiraTarget.Do).
			AddTask(telegramSource.Do).
			Start(mainCtx)
		pkg.Throw(err)
	}()
	select {}
}

func parseArgs(args []string) []string {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("config", "./cmd/config.yaml", "path to config")

	err := flags.Parse(args[1:])
	pkg.Throw(err)

	return []string{*configPath}
}
