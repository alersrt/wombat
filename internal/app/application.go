package app

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"wombat/internal/dao"
	"wombat/internal/domain"
	"wombat/pkg/daemon"
	"wombat/pkg/errors"
)

type MessageHelper interface {
	SendToTopic(topic string, key []byte, message []byte) error
	Subscribe(topics []string, handler func(*kafka.Message) error) error
}

type Source interface {
	ForwardTo(chan *domain.MessageEvent)
}

type Application struct {
	executor               *daemon.Daemon
	conf                   *Config
	sourceChan             chan *domain.MessageEvent
	kafkaHelper            MessageHelper
	messageEventRepository dao.QueryHelper[domain.MessageEvent, string]
	telegramSource         Source
}

func NewApplication(
	executor *daemon.Daemon,
	kafkaHelper MessageHelper,
	messageEventRepository dao.QueryHelper[domain.MessageEvent, string],
	telegramSource Source,
) (*Application, error) {
	conf, ok := executor.GetConfig().(*Config)
	if !ok {
		return nil, errors.NewError("Wrong config type")
	}
	return &Application{
		conf:                   conf,
		executor:               executor,
		sourceChan:             make(chan *domain.MessageEvent),
		kafkaHelper:            kafkaHelper,
		messageEventRepository: messageEventRepository,
		telegramSource:         telegramSource,
	}, nil
}

func (receiver *Application) Run(ctx context.Context) {
	go receiver.executor.
		AddTask(receiver.send).
		AddTask(receiver.route).
		AddTask(receiver.source).
		Start(ctx)

	select {}
}
