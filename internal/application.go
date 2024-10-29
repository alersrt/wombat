package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"regexp"
	"strconv"
	"wombat/internal/config"
	"wombat/internal/domain"
	"wombat/internal/messaging"
	"wombat/internal/source"
	"wombat/pkg/daemon"
)

type Application interface {
	Run()
}

type application struct {
	mainCtx             context.Context
	mainCancelCauseFunc context.CancelCauseFunc
	updates             chan any
	conf                *config.Config
	kafkaHelper         messaging.KafkaHelper
	telegram            source.Source
}

func NewApplication(
	mainCtx context.Context,
	mainCancelCauseFunc context.CancelCauseFunc,
	updates chan any,
	conf *config.Config,
	kafkaHelper messaging.KafkaHelper,
	telegram source.Source,
) Application {
	return &application{
		mainCtx:             mainCtx,
		mainCancelCauseFunc: mainCancelCauseFunc,
		conf:                conf,
		updates:             updates,
		kafkaHelper:         kafkaHelper,
		telegram:            telegram,
	}
}

func (receiver *application) Run() {
	dmn := daemon.Create(receiver.mainCtx, receiver.mainCancelCauseFunc, receiver.conf)

	go dmn.Start(receiver.readFromTopic)
	go dmn.Start(receiver.readUpdates)
	go dmn.Start(receiver.telegram.Read)

	select {}
}

func (receiver *application) readUpdates(cancel context.CancelCauseFunc) {
	for update := range receiver.updates {
		switch matched := update.(type) {

		case tgbotapi.Update:
			if matched.Message != nil {
				pattern := regexp.MustCompile(receiver.conf.Bot.Tag)
				tags := pattern.FindAllString(matched.Message.Text, -1)

				if len(tags) > 0 {

					msg := &domain.MessageEvent{
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

func (receiver *application) readFromTopic(cancel context.CancelCauseFunc) {
	err := receiver.kafkaHelper.Subscribe([]string{receiver.conf.Kafka.Topic}, func(event *kafka.Message) error {
		msg := &domain.MessageEvent{}
		if err := json.Unmarshal(event.Value, msg); err != nil {
			slog.Warn(err.Error())
			return err
		}

		slog.Info(fmt.Sprintf(
			"Consumed event from topic %s: key = %-10s value = %+v",
			event.TopicPartition.Topic,
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
