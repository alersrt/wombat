package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"wombat/internal/config"
	"wombat/internal/jira"
	"wombat/internal/router"
	"wombat/internal/storage"
	"wombat/internal/telegram"
	"wombat/pkg/cipher"
)

var (
	ErrApp     = errors.New("app")
	ErrAppNew  = errors.New("app: new")
	ErrAppInit = errors.New("app: init")
)

type App struct {
	mtx    sync.Mutex
	source *telegram.Source
	target *jira.Target
}

func (a *App) Init(args []string) error {
	slog.Info("app: init: start")
	defer slog.Info("app: init: finish")
	a.mtx.Lock()
	defer a.mtx.Unlock()

	conf := new(config.Config)
	err := conf.Init(args)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAppInit, err)
	}

	gcm, err := cipher.NewAesGcmCipher([]byte(conf.Cipher.Key))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAppInit, err)
	}
	rt := router.NewRouter()

	db, err := storage.NewDbStorage(conf.PostgreSQL.Url)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAppInit, err)
	}
	tgTarget, err := telegram.NewTelegramSource(conf.Telegram.Token, rt, db, gcm)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAppInit, err)
	}
	jiraTarget := jira.NewJiraTarget(conf.Jira.Url, conf.Bot.Tag, rt, db, gcm)

	a.source = tgTarget
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
