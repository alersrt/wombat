package internal

import (
	"context"
	"log/slog"
	"sync"
	"wombat/internal/jira"
	"wombat/internal/telegram"
)

type App struct {
	mtx    sync.Mutex
	source *telegram.Source
	target *jira.Target
}

func (a *App) Init() error {

	return nil
}

func (a *App) Do(ctx context.Context) {
	slog.Info("app: do: start")
	defer slog.Info("app: do: finish")

	go a.source.DoReq(ctx)
	go a.source.DoRes(ctx)
	go a.target.Do(ctx)

	<-ctx.Done()
}
