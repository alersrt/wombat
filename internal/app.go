package internal

import (
	"context"
	"log/slog"
	"sync"
	"wombat/internal/config"
	"wombat/internal/jira"
	router2 "wombat/internal/router"
	"wombat/internal/storage"
	"wombat/internal/telegram"
	"wombat/pkg/cipher"
)

type App struct {
	mtx    sync.Mutex
	source *telegram.Source
	target *jira.Target
}

func (a *App) Init(args []string) error {
	slog.Info("app:init:start")
	defer slog.Info("app:init:finish")
	a.mtx.Lock()
	defer a.mtx.Unlock()

	conf := new(config.Config)
	err := conf.Init(args)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewAesGcmCipher([]byte(conf.Cipher.Key))
	if err != nil {
		return err
	}
	router := router2.NewRouter()

	db, err := storage.NewDbStorage(conf.PostgreSQL.Url)
	if err != nil {
		return err
	}
	telegramSource, err := telegram.NewTelegramSource(conf.Telegram.Token, router, db, gcm)
	if err != nil {
		return err
	}
	jiraTarget := jira.NewJiraTarget(conf.Jira.Url, conf.Bot.Tag, router, db, gcm)

	a.source = telegramSource
	a.target = jiraTarget

	return nil
}

func (a *App) Do(ctx context.Context) {
	slog.Info("app:do:start")
	defer slog.Info("app:do:finish")

	go a.source.DoReq(ctx)
	go a.source.DoRes(ctx)
	go a.target.Do(ctx)

	<-ctx.Done()
}
