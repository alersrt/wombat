package app

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log/slog"
	"wombat/internal/domain"
)

func (receiver *Application) send() {
	err := receiver.kafkaHelper.Subscribe([]string{receiver.conf.Kafka.Topic}, func(event *kafka.Message) error {
		msg := &domain.MessageEvent{}
		if err := json.Unmarshal(event.Value, msg); err != nil {
			return err
		}

		saved, err := receiver.messageEventRepository.Save(msg)

		if err != nil {
			return err
		}
		slog.Info(fmt.Sprintf("Consume message: %s", string(event.Key)))
		slog.Info(fmt.Sprintf("Value: %+v", saved))

		return nil
	})

	if err != nil {
		slog.Warn(err.Error())
	}
}
