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
	SendToTopic(topic string, message []byte) error
	Subscribe(topics []string, handler func(*kafka.Message) error) error
}

type Source interface {
	ForwardTo(chan any)
}

type Application struct {
	executor               *daemon.Daemon
	conf                   *Config
	routeChan              chan any
	kafkaHelper            MessageHelper
	messageEventRepository dao.QueryHelper[domain.MessageEvent, string]
	telegram               Source
}

func NewApplication(
	executor *daemon.Daemon,
	routeChan chan any,
	kafkaHelper MessageHelper,
	messageEventRepository dao.QueryHelper[domain.MessageEvent, string],
	telegram Source,
) (*Application, error) {
	conf, ok := executor.GetConfig().(*Config)
	if !ok {
		return nil, errors.NewError("Wrong config type")
	}
	return &Application{
		conf:                   conf,
		executor:               executor,
		routeChan:              routeChan,
		kafkaHelper:            kafkaHelper,
		messageEventRepository: messageEventRepository,
		telegram:               telegram,
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
