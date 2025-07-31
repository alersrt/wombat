package app

import (
	"context"
	"wombat/internal/domain"
	"wombat/pkg/daemon"
	"wombat/pkg/errors"
)

type Source interface {
	GetSourceType() domain.SourceType
	Process()
}

type Target interface {
	GetTargetType() domain.TargetType
	Process()
}

type Application struct {
	executor       *daemon.Daemon
	conf           *Config
	forwardChannel chan *domain.Message
	source         Source
	target         Target
}

func NewApplication(
	executor *daemon.Daemon,
	source Source,
	target Target,
) (*Application, error) {
	conf, ok := executor.GetConfig().(*Config)
	if !ok {
		return nil, errors.NewError("Wrong config type")
	}
	return &Application{
		conf:           conf,
		executor:       executor,
		forwardChannel: make(chan *domain.Message),
		source:         source,
		target:         target,
	}, nil
}

func (receiver *Application) Run(ctx context.Context) {
	go receiver.executor.
		AddTask(receiver.target.Process).
		AddTask(receiver.source.Process).
		Start(ctx)

	select {}
}
