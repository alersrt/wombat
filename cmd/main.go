package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"regexp"
	"time"
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
		log.ErrorLog.Fatal(err)
	}

	daemon.Create(conf, func(cancel context.CancelCauseFunc) {
		bot.Debug = true

		log.InfoLog.Print("Authorized on account %s", bot.Self.UserName)

		for update := range getUpdatesFromBot(bot) {

			if update.EditedMessage != nil {
				pattern := regexp.MustCompile(conf.Bot.Tag)
				tag := pattern.FindAllString(update.EditedMessage.Text, -1)
				key := fmt.Sprintf(
					"%d-%d",
					update.EditedMessage.Chat.ID,
					update.EditedMessage.MessageID,
				)
				err := sendToTopic(producer, conf.Kafka.Topic, []byte(key), []byte(update.EditedMessage.Text))
				if err != nil {
					log.WarningLog.Print(err)
				}
				log.InfoLog.Printf("Sent message: %s => %s", tag, key)
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
		log.ErrorLog.Fatal(err)
	}
	return producer
}

func getKafkaConsumer(configMap *kafka.ConfigMap) *kafka.Consumer {
	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		log.ErrorLog.Fatal(err)
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
		log.ErrorLog.Fatal(err)
	}

	// Process messages
	for {
		ev, err := consumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			//log.WarningLog.Print(err)
			continue
		}

		log.InfoLog.Printf(
			"Consumed event from topic %s: key = %-10s value = %s\n",
			*ev.TopicPartition.Topic,
			string(ev.Key),
			string(ev.Value),
		)
	}

	consumer.Close()
}
