package app

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"regexp"
	"wombat/internal/domain"
	"wombat/internal/storage"
	"wombat/pkg/daemon"
	"wombat/pkg/errors"
)

type MessageHelper interface {
	SendToTopic(topic string, key []byte, message []byte) error
	Subscribe(topics []string, handler func(*kafka.Message) error) error
}

type TargetClient interface {
	Add(tag string, text string) (string, error)
	Update(tag string, commentId string, text string) error
}

type Source interface {
	ForwardTo(chan *domain.Message)
}

type Application struct {
	executor             *daemon.Daemon
	conf                 *Config
	sourceChan           chan *domain.Message
	kafkaHelper          MessageHelper
	aclRepository        *storage.AclRepository
	commentRepository    *storage.CommentRepository
	connectionRepository *storage.ConnectionRepository
	targetClientFactory  func(targetType domain.TargetType, sourceType domain.SourceType, authorId string) (TargetClient, error)
	telegramSource       Source
	tagsRegex            *regexp.Regexp
}

func NewApplication(
	executor *daemon.Daemon,
	kafkaHelper MessageHelper,
	aclRepository *storage.AclRepository,
	commentRepository *storage.CommentRepository,
	connectionRepository *storage.ConnectionRepository,
	targetClientFactory func(targetType domain.TargetType, sourceType domain.SourceType, authorId string) (TargetClient, error),
	telegramSource Source,
) (*Application, error) {
	conf, ok := executor.GetConfig().(*Config)
	if !ok {
		return nil, errors.NewError("Wrong config type")
	}
	return &Application{
		conf:                 conf,
		executor:             executor,
		sourceChan:           make(chan *domain.Message),
		kafkaHelper:          kafkaHelper,
		aclRepository:        aclRepository,
		commentRepository:    commentRepository,
		connectionRepository: connectionRepository,
		targetClientFactory:  targetClientFactory,
		telegramSource:       telegramSource,
		tagsRegex:            regexp.MustCompile(conf.Bot.Tag),
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
