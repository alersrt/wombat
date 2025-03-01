package app

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"regexp"
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
	ForwardTo(chan *domain.Message)
}

type Application struct {
	executor          *daemon.Daemon
	conf              *Config
	sourceChan        chan *domain.Message
	kafkaHelper       MessageHelper
	jiraHelper        JiraHelper
	aclRepository     *dao.AclRepository
	commentRepository *dao.CommentRepository
	telegramSource    Source
	tagsRegex         *regexp.Regexp
}

func NewApplication(
	executor *daemon.Daemon,
	kafkaHelper MessageHelper,
	jiraHelper JiraHelper,
	aclRepository *dao.AclRepository,
	commentRepository *dao.CommentRepository,
	telegramSource Source,
) (*Application, error) {
	conf, ok := executor.GetConfig().(*Config)
	if !ok {
		return nil, errors.NewError("Wrong config type")
	}
	return &Application{
		conf:              conf,
		executor:          executor,
		sourceChan:        make(chan *domain.Message),
		kafkaHelper:       kafkaHelper,
		jiraHelper:        jiraHelper,
		aclRepository:     aclRepository,
		commentRepository: commentRepository,
		telegramSource:    telegramSource,
		tagsRegex:         regexp.MustCompile(conf.Bot.Tag),
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
