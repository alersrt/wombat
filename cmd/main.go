package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"wombat/internal/config"
	"wombat/internal/message"
	"wombat/pkg/daemon"
)

var (
	mainCtx             context.Context
	mainCancelCauseFunc context.CancelCauseFunc
	conf                *config.Config
	producer            *kafka.Producer
	consumer            *kafka.Consumer
	updates             chan any
)

func main() {
	mainCtx, mainCancelCauseFunc = context.WithCancelCause(context.Background())

	conf = new(config.Config)
	err := conf.Init(os.Args)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers": conf.Kafka.Bootstrap,
		"group.id":          conf.Kafka.GroupId,
		"auto.offset.reset": "earliest",
	}
	producer = getKafkaProducer(kafkaConf)
	consumer = getKafkaConsumer(kafkaConf)

	updates = make(chan any)

	dmn := daemon.Create(mainCtx, mainCancelCauseFunc, conf)
	go dmn.Start(readFromTopic)
	go dmn.Start(readUpdates)

	telegram := &Telegram{Updates: updates}
	go dmn.Start(telegram.Read)

	select {}
}

func readUpdates(cancel context.CancelCauseFunc) {
	for update := range updates {
		switch matched := update.(type) {

		case tgbotapi.Update:
			if matched.Message != nil {
				pattern := regexp.MustCompile(conf.Bot.Tag)
				tags := pattern.FindAllString(matched.Message.Text, -1)

				if len(tags) > 0 {

					msg := &message.MessageEvent{
						SourceType: message.TELEGRAM,
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

					err = sendToTopic(conf.Kafka.Topic, jsonifiedMsg)
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

func readFromTopic(cancel context.CancelCauseFunc) {
	err := consumer.SubscribeTopics([]string{conf.Kafka.Topic}, nil)
	if err != nil {
		slog.Error(err.Error())
		cancel(err)
	}

	// Process messages
	for {
		ev, err := consumer.ReadMessage(-1)
		if err != nil {
			slog.Warn(err.Error())
			continue
		}

		msg := &message.MessageEvent{}
		if err = json.Unmarshal(ev.Value, msg); err != nil {
			slog.Warn(err.Error())
			continue
		}

		slog.Info(fmt.Sprintf(
			"Consumed event from topic %s: key = %-10s value = %+v",
			*ev.TopicPartition.Topic,
			string(ev.Key),
			msg,
		))
	}

	consumer.Close()
}

func getKafkaProducer(configMap *kafka.ConfigMap) *kafka.Producer {
	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return producer
}

func getKafkaConsumer(configMap *kafka.ConfigMap) *kafka.Consumer {
	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return consumer
}

func sendToTopic(topic string, message []byte) error {
	key, err := uuid.New().MarshalText()
	if err != nil {
		slog.Warn(err.Error())
		return err
	}
	return producer.Produce(&kafka.Message{
		Key:   key,
		Value: message,
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
	}, nil)
}
