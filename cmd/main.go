package main

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"os"
	"regexp"
	"wombat/internal/config"
	"wombat/pkg/daemon"
	"wombat/pkg/log"
)

func main() {
	conf := &config.Config{}
	err := conf.Init(os.Args)
	if err != nil {
		log.ErrorLog.Fatal(err)
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": conf.Kafka.Bootstrap,
		"client.id":         conf.Kafka.ClientId,
		"acks":              "all"})
	if err != nil {
		log.ErrorLog.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(conf.Telegram.Token)
	if err != nil {
		log.ErrorLog.Fatal(err)
	}

	daemon.Create(conf, func(cancel context.CancelCauseFunc) {
		bot.Debug = true

		log.InfoLog.Print("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.AllowedUpdates = append(u.AllowedUpdates, tgbotapi.UpdateTypeMessageReaction, tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeEditedMessage)
		u.Timeout = 60

		updates := bot.GetUpdatesChan(u)

		for update := range updates {

			if update.EditedMessage != nil {
				pattern := regexp.MustCompile(conf.Bot.Tag)
				tag := pattern.FindAllString(update.EditedMessage.Text, -1)
				newUuid, _ := uuid.New().MarshalBinary()
				err := producer.Produce(&kafka.Message{
					Key:   newUuid,
					Value: []byte(update.EditedMessage.Text),
					TopicPartition: kafka.TopicPartition{
						Topic:     &conf.Kafka.Topic,
						Partition: kafka.PartitionAny,
					},
				}, nil)
				if err != nil {
					log.ErrorLog.Print(err)
				}
				log.InfoLog.Print(tag)
				log.InfoLog.Print(update.EditedMessage.Text)
			}
		}
	})
}
