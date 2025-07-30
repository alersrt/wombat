package app

import (
	"context"
	"regexp"
	"wombat/internal/domain"
	"wombat/internal/storage"
	"wombat/pkg/daemon"
	"wombat/pkg/errors"
)

type Target interface {
	Add(tag string, text string) (string, error)
	Update(tag string, commentId string, text string) error
}

type Source interface {
	ForwardTo(chan *domain.Message)
}

type Application struct {
	executor   *daemon.Daemon
	conf       *Config
	sourceChan chan *domain.Message
	targetChan chan *domain.Message
	dbStorage  *storage.DbStorage
	source     Source
	tagsRegex  *regexp.Regexp
}

func NewApplication(
	executor *daemon.Daemon,
	dbStorage *storage.DbStorage,
	source Source,
) (*Application, error) {
	conf, ok := executor.GetConfig().(*Config)
	if !ok {
		return nil, errors.NewError("Wrong config type")
	}
	return &Application{
		conf:       conf,
		executor:   executor,
		sourceChan: make(chan *domain.Message),
		targetChan: make(chan *domain.Message),
		dbStorage:  dbStorage,
		source:     source,
		tagsRegex:  regexp.MustCompile(conf.Bot.Tag),
	}, nil
}

func (receiver *Application) Run(ctx context.Context) {
	go receiver.executor.
		AddTask(receiver.processTarget).
		AddTask(receiver.route).
		AddTask(receiver.processSource).
		Start(ctx)

	select {}
}
