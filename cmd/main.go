package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"regexp"
	"wombat/internal/config"
	"wombat/pkg/daemon"
	"wombat/pkg/log"
)

func main() {
	conf := new(config.Config)
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

		for update := range getUpdatesChannel(bot) {

			if update.EditedMessage != nil {
				pattern := regexp.MustCompile(conf.Bot.Tag)
				tag := pattern.FindAllString(update.EditedMessage.Text, -1)
				key := fmt.Sprintf(
					"%d-%d",
					update.EditedMessage.Chat.ID,
					update.EditedMessage.MessageID,
				)
				err := SendTo(producer, conf.Kafka.Topic, []byte(key), []byte(update.EditedMessage.Text))
				if err != nil {
					log.WarningLog.Print(err)
				}
				log.InfoLog.Printf("Send: %s => %s", tag, key)
			}
		}
	})
}

func getUpdatesChannel(api *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
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

func SendTo(producer *kafka.Producer, topic string, key []byte, message []byte) error {
	return producer.Produce(&kafka.Message{
		Key:   key,
		Value: message,
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
	}, nil)
}
