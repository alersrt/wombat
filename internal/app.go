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
	mu     sync.Mutex
	wg     sync.WaitGroup
	cfg    *config.Config
	source *TelegramSource
	target *JiraTarget
}

func (a *App) Init(args []string) error {
	slog.Info("app:init:start")
	defer slog.Info("app:init:finish")
	a.mu.Lock()
	defer a.mu.Unlock()

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

	// Any of these three functions are important for execution so the goroutine must
	// be stopped if any of these is stopped.
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		go a.source.DoReq(ctx)
	}()
	go func() {
		defer a.wg.Done()
		go a.source.DoRes(ctx)
	}()
	go func() {
		defer a.wg.Done()
		go a.target.Do(ctx)
	}()
	a.wg.Wait()
}

func (a *App) Shutdown() {
	slog.Info("app:shutdown:start")
	defer slog.Info("app:shutdown:finish")
	a.wg.Add(-3)
}
