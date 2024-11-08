package app

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
)

func (receiver *Application) route() {
	for update := range receiver.sourceChan {
		if update.AuthorId != "" && receiver.tagsRegex.MatchString(update.Text) {
			jsonifiedMsg, err := json.Marshal(update)
			if err != nil {
				slog.Warn(err.Error())
				return
			}

			key, err := uuid.New().MarshalBinary()
			if err != nil {
				slog.Warn(err.Error())
				return
			}
			err = receiver.kafkaHelper.SendToTopic(receiver.conf.Kafka.Topic, key, jsonifiedMsg)
			if err != nil {
				slog.Warn(err.Error())
				return
			}
			slog.Info(fmt.Sprintf("Send message: %s", string(key)))
		}
	}
}
