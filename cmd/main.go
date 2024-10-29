package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"wombat/internal/config"
	"wombat/internal/message"
	"wombat/pkg/daemon"
)

var (
	ctx      context.Context
	cancel   context.CancelCauseFunc
	bot      *tgbotapi.BotAPI
	conf     *config.Config
	producer *kafka.Producer
	consumer *kafka.Consumer
)

func main() {
	ctx, cancel = context.WithCancelCause(context.Background())

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

	bot, err = tgbotapi.NewBotAPI(conf.Telegram.Token)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	dmn := daemon.Create(ctx, cancel, conf)
	go dmn.Start(readFromTopic)
	go dmn.Start(readFromBot)

	select {}
}

func readFromBot(cancel context.CancelCauseFunc) {
	slog.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	for update := range getUpdatesFromBot(bot) {

		if update.Message != nil {
			pattern := regexp.MustCompile(conf.Bot.Tag)
			tags := pattern.FindAllString(update.Message.Text, -1)

			if len(tags) > 0 {
				key := fmt.Sprintf(
					"%d-%d",
					update.Message.Chat.ID,
					update.Message.MessageID,
				)

				msg := &message.MessageEvent{
					SourceType: message.TELEGRAM,
					Text:       update.Message.Text,
					AuthorId:   update.Message.From.UserName,
					ChatId:     strconv.FormatInt(update.Message.Chat.ID, 10),
					MessageId:  strconv.Itoa(update.Message.MessageID),
				}

				jsonifiedMsg, err := json.Marshal(msg)
				if err != nil {
					slog.Warn(err.Error())
					return
				}

				err = sendToTopic(conf.Kafka.Topic, []byte(key), jsonifiedMsg)
				if err != nil {
					slog.Warn(err.Error())
					return
				}
				slog.Info(fmt.Sprintf("Sent message: %s => %s", tags, key))
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
			"Consumed event from topic %s: key = %-10s value = %s",
			*ev.TopicPartition.Topic,
			string(ev.Key),
			msg,
		))
	}

	consumer.Close()
}

func getUpdatesFromBot(api *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.AllowedUpdates = append(
		u.AllowedUpdates,
		tgbotapi.UpdateTypeMessageReaction,
		tgbotapi.UpdateTypeMessage,
		tgbotapi.UpdateTypeEditedMessage,
	)
	u.Timeout = 60
	return api.GetUpdatesChan(u)
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

func sendToTopic(topic string, key []byte, message []byte) error {
	return producer.Produce(&kafka.Message{
		Key:   key,
		Value: message,
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
	}, nil)
}
