package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"log/slog"
	"os"
	"wombat/internal"
	"wombat/pkg"
)

func main() {
	mainCtx, mainCancelCauseFunc := context.WithCancelCause(context.Background())
	defer mainCancelCauseFunc(nil)

	conf := new(internal.Config)
	args, err := parseArgs(os.Args)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

	err = conf.Init(args)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

	cipher, err := internal.NewAesGcmCipher([]byte(conf.Cipher.Key))
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
	router := internal.NewRouter()

	db, err := internal.NewDbStorage(conf.PostgreSQL.Url)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
	telegramSource, err := internal.NewTelegramSource(conf.Telegram.Token, router, db, cipher)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
	jiraTarget := internal.NewJiraTarget(conf.Jira.Url, conf.Bot.Tag, router, db, cipher)

	dmn := pkg.Create(conf)
	go func() {
		code, err := dmn.
			AddTask(jiraTarget).
			AddTask(telegramSource).
			Start(mainCtx)
		if err != nil {
			slog.Error(fmt.Sprintf("%+v", err))
			os.Exit(code)
		}
	}()
	select {}
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
