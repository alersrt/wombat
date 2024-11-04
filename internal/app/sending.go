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

		saved, err := receiver.messageEventRepository.GetById(msg.Hash)
		if err != nil {
			return err
		}

		if saved == nil || saved.CommentId == "" {
			for _, tag := range msg.Tags {
				commentId, err := receiver.jiraHelper.AddComment(tag, msg.Text)
				if err != nil {
					return err
				}
				msg.CommentId = commentId
				_, err = receiver.messageEventRepository.Save(msg)
				if err != nil {
					return err
				}
			}
		} else {
			for _, tag := range msg.Tags {
				err = receiver.jiraHelper.UpdateComment(tag, saved.CommentId, msg.Text)
				if err != nil {
					return err
				}
				_, err = receiver.messageEventRepository.Save(msg)
				if err != nil {
					return err
				}
			}
		}

		slog.Info(fmt.Sprintf("Consume message: %s", string(event.Key)))
		slog.Info(fmt.Sprintf("Value: %+v", saved))

		return nil
	})

	if err != nil {
		slog.Warn(err.Error())
	}
}
