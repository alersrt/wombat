package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"log/slog"
	"os"
	"sync"
	"wombat/internal"
	"wombat/pkg"
)

type App struct {
	mtx    sync.Mutex
	wg     sync.WaitGroup
	conf   *pkg.Config
	source *internal.TelegramSource
	target *internal.JiraTarget
}

var _ pkg.Daemon = (*App)(nil)

func (a *App) Init(args []string) error {
	slog.Info("app:init:start")
	defer slog.Info("app:init:finish")
	a.mtx.Lock()
	defer a.mtx.Unlock()

	conf := new(internal.Config)
	err := conf.Init(args)
	if err != nil {
		return err
	}

	cipher, err := internal.NewAesGcmCipher([]byte(conf.Cipher.Key))
	if err != nil {
		return err
	}
	router := internal.NewRouter()

	db, err := internal.NewDbStorage(conf.PostgreSQL.Url)
	if err != nil {
		return err
	}
	telegramSource, err := internal.NewTelegramSource(conf.Telegram.Token, router, db, cipher)
	if err != nil {
		return err
	}
	jiraTarget := internal.NewJiraTarget(conf.Jira.Url, conf.Bot.Tag, router, db, cipher)

	a.source = telegramSource
	a.target = jiraTarget

	return nil
}

func (a *App) Do(ctx context.Context) {
	slog.Info("app:do:start")
	defer slog.Info("app:do:finish")
	a.wg.Add(1)
	go a.source.Do(ctx, &a.wg)
	a.wg.Add(1)
	go a.target.Do(ctx, &a.wg)
	a.wg.Wait()
}

func (a *App) Shutdown() {
	slog.Info("app:shutdown:start")
	defer slog.Info("app:shutdown:finish")
	a.wg.Add(-2)
}

func main() {
	ctx := context.Background()

	app := new(App)
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

	code, err := pkg.HandleSignals(ctx, app)
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
