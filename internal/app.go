package internal

import (
	"context"
	"log/slog"
	"sync"
	"wombat/internal/config"
	"wombat/internal/storage"
	"wombat/pkg/cipher"
)

type App struct {
	mtx    sync.Mutex
	cfg    *config.Config
	source *TelegramSource
	target *JiraTarget
}

func (a *App) Init(ctx context.Context, args []string) error {
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
	router := NewRouter()

	db, err := storage.NewDbStorage(conf.PostgreSQL.Url)
	if err != nil {
		return err
	}
	telegramSource, err := NewTelegramSource(conf.Telegram.Token, router, db, gcm)
	if err != nil {
		return err
	}
	jiraTarget := NewJiraTarget(conf.Jira.Url, conf.Bot.Tag, router, db, gcm)

	a.source = telegramSource
	a.target = jiraTarget

	return nil
}

func (a *App) Do(ctx context.Context) {
	slog.Info("app:do:start")
	defer slog.Info("app:do:finish")

	go func() {
		go a.source.DoReq(ctx)
	}()
	go func() {
		go a.source.DoRes(ctx)
	}()
	go func() {
		go a.target.Do(ctx)
	}()

}

func (a *App) Shutdown(cancel func()) {
	slog.Info("app:shutdown:start")
	defer slog.Info("app:shutdown:finish")
	cancel()
}
