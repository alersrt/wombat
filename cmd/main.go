package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
	"regexp"
	"wombat/internal/config"
	"wombat/pkg/daemon"
)

func main() {
	conf := new(config.Config)
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
	producer := getKafkaProducer(kafkaConf)
	consumer := getKafkaConsumer(kafkaConf)

	go readFromTopic(consumer, conf.Kafka.Topic)

	bot, err := tgbotapi.NewBotAPI(conf.Telegram.Token)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	daemon.Create(conf, func(cancel context.CancelCauseFunc) {
		slog.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

		for update := range getUpdatesFromBot(bot) {

			if update.EditedMessage != nil {
				pattern := regexp.MustCompile(conf.Bot.Tag)
				tags := pattern.FindAllString(update.EditedMessage.Text, -1)

				if len(tags) > 0 {
					key := fmt.Sprintf(
						"%d-%d",
						update.EditedMessage.Chat.ID,
						update.EditedMessage.MessageID,
					)
					err := sendToTopic(producer, conf.Kafka.Topic, []byte(key), []byte(update.EditedMessage.Text))
					if err != nil {
						slog.Warn(err.Error())
					}
					slog.Info(fmt.Sprintf("Sent message: %s => %s", tags, key))
				}
			}
		}
	})
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

func sendToTopic(producer *kafka.Producer, topic string, key []byte, message []byte) error {
	return producer.Produce(&kafka.Message{
		Key:   key,
		Value: message,
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
	}, nil)
}

func readFromTopic(consumer *kafka.Consumer, topic string) {
	err := consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// Process messages
	for {
		ev, err := consumer.ReadMessage(-1)
		if err != nil {
			slog.Warn(err.Error())
			continue
		}

		slog.Info(fmt.Sprintf(
			"Consumed event from topic %s: key = %-10s value = %s",
			*ev.TopicPartition.Topic,
			string(ev.Key),
			string(ev.Value),
		))
	}

	consumer.Close()
}
