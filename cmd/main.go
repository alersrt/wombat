package main

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"regexp"
	"wombat/internal/config"
	"wombat/pkg/daemon"
	"wombat/pkg/log"
)

func main() {
	conf := &config.Config{}

	daemon.Create(conf, func(cancel context.CancelCauseFunc) {
		producer, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": conf.Kafka.Bootstrap,
			"client.id":         conf.Kafka.ClientId,
			"acks":              "all"})
		if err != nil {
			log.Error(err)
			cancel(err)
		}
		bot, err := tgbotapi.NewBotAPI(conf.Telegram.Token)
		if err != nil {
			log.Error(err)
			cancel(err)
		}

		bot.Debug = true

		log.Infof("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.AllowedUpdates = append(u.AllowedUpdates, tgbotapi.UpdateTypeMessageReaction, tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeEditedMessage)
		u.Timeout = 60

		updates := bot.GetUpdatesChan(u)

		for update := range updates {

			if update.EditedMessage != nil {
				pattern := regexp.MustCompile(conf.Bot.Tag)
				tag := pattern.FindAllString(update.EditedMessage.Text, -1)
				newUuid, _ := uuid.New().MarshalBinary()
				topic := "topic"
				producer.Produce(&kafka.Message{
					Key:   newUuid,
					Value: []byte(update.EditedMessage.Text),
					TopicPartition: kafka.TopicPartition{
						Topic:     &topic,
						Partition: kafka.PartitionAny,
					},
				}, nil)
				log.Info(tag)
				log.Info(update.EditedMessage.Text)
			}
		}
	})
}
