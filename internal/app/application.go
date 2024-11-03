package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"log/slog"
	"regexp"
	"strconv"
	"wombat/internal/config"
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
	conf                   *config.Config
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
	conf, ok := executor.GetConfig().(*config.Config)
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
		AddTask(receiver.route).
		AddTask(receiver.source).
		AddTask(receiver.forwardFromTelegram).
		Start(ctx)

	select {}
}

func (receiver *Application) forwardFromTelegram() {
	receiver.telegram.ForwardTo(receiver.routeChan)
}

func (receiver *Application) source() {

	hash := func(sourceType domain.SourceType, chatId string, messageId string) string {
		return uuid.NewSHA1(uuid.NameSpaceURL, []byte(sourceType.String()+chatId+messageId)).String()
	}

	for update := range receiver.routeChan {
		switch matched := update.(type) {
		case tgbotapi.Update:
			if matched.Message != nil {
				pattern := regexp.MustCompile(receiver.conf.Bot.Tag)
				tags := pattern.FindAllString(matched.Message.Text, -1)

				if len(tags) > 0 {
					messageId := strconv.Itoa(matched.Message.MessageID)
					chatId := strconv.FormatInt(matched.Message.Chat.ID, 10)

					msg := &domain.MessageEvent{
						Hash:       hash(domain.TELEGRAM, chatId, messageId),
						EventType:  domain.CREATE,
						SourceType: domain.TELEGRAM,
						Text:       matched.Message.Text,
						AuthorId:   matched.Message.From.UserName,
						ChatId:     chatId,
						MessageId:  messageId,
					}

					jsonifiedMsg, err := json.Marshal(msg)
					if err != nil {
						slog.Warn(err.Error())
						return
					}

					err = receiver.kafkaHelper.SendToTopic(receiver.conf.Kafka.Topic, jsonifiedMsg)
					if err != nil {
						slog.Warn(err.Error())
						return
					}
					slog.Info(fmt.Sprintf("Sent message: %s", tags))
				}
			}
		}
	}
}

func (receiver *Application) route() {
	err := receiver.kafkaHelper.Subscribe([]string{receiver.conf.Kafka.Topic}, func(event *kafka.Message) error {
		msg := &domain.MessageEvent{}
		if err := json.Unmarshal(event.Value, msg); err != nil {
			return err
		}

		saved, err := receiver.messageEventRepository.Save(msg)
		if err != nil {
			return err
		}
		slog.Info(fmt.Sprintf(
			"Consumed event from topic %s: key = %-10s value = %+v",
			*event.TopicPartition.Topic,
			string(event.Key),
			saved,
		))
		return nil
	})

	if err != nil {
		slog.Warn(err.Error())
	}
}
