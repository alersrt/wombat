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
	"wombat/internal/domain"
	"wombat/internal/messaging"
	"wombat/internal/source"
	"wombat/pkg/daemon"
	"wombat/pkg/errors"
)

type Application interface {
	Run()
}

type application struct {
	executor            *daemon.Daemon
	conf                *config.Config
	mainCtx             context.Context
	mainCancelCauseFunc context.CancelCauseFunc
	updates             chan any
	kafkaHelper         messaging.KafkaHelper
	telegram            source.Source
}

func NewApplication(
	executor *daemon.Daemon,
	updates chan any,
	kafkaHelper messaging.KafkaHelper,
	telegram source.Source,
) (Application, error) {
	conf, ok := executor.GetConfig().(*config.Config)
	if !ok {
		return nil, errors.NewError("Wrong config type")
	}
	return &application{
		mainCtx:             executor.GetContext(),
		mainCancelCauseFunc: executor.GetCancelCauseFunc(),
		conf:                conf,
		executor:            executor,
		updates:             updates,
		kafkaHelper:         kafkaHelper,
		telegram:            telegram,
	}, nil
}

func (receiver *application) Run() {
	go receiver.executor.Start(receiver.route)
	go receiver.executor.Start(receiver.source)
	go receiver.executor.Start(receiver.telegram.Read)

	select {}
}

func (receiver *application) source(cancel context.CancelCauseFunc) {
	for update := range receiver.updates {
		switch matched := update.(type) {

		case tgbotapi.Update:
			if matched.Message != nil {
				pattern := regexp.MustCompile(receiver.conf.Bot.Tag)
				tags := pattern.FindAllString(matched.Message.Text, -1)

				if len(tags) > 0 {

					msg := &domain.MessageEvent{
						EventType:  domain.CREATE,
						SourceType: domain.TELEGRAM,
						Text:       matched.Message.Text,
						AuthorId:   matched.Message.From.UserName,
						ChatId:     strconv.FormatInt(matched.Message.Chat.ID, 10),
						MessageId:  strconv.Itoa(matched.Message.MessageID),
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

func (receiver *application) route(cancel context.CancelCauseFunc) {
	err := receiver.kafkaHelper.Subscribe([]string{receiver.conf.Kafka.Topic}, func(event *kafka.Message) error {
		msg := &domain.MessageEvent{}
		if err := json.Unmarshal(event.Value, msg); err != nil {
			slog.Warn(err.Error())
			return err
		}

		hash := uuid.NewSHA1(uuid.NameSpaceURL, []byte(msg.ChatId+msg.MessageId))

		slog.Info(hash.String())

		slog.Info(fmt.Sprintf(
			"Consumed event from topic %s: key = %-10s value = %+v",
			*event.TopicPartition.Topic,
			string(event.Key),
			msg,
		))
		return nil
	})

	if err != nil {
		slog.Error(err.Error())
		cancel(err)
	}
}
