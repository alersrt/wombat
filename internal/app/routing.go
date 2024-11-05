package app

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
)

func (receiver *Application) route() {
	for update := range receiver.sourceChan {
		pattern := regexp.MustCompile(receiver.conf.Bot.Tag)
		tags := pattern.FindAllString(update.Text, -1)

		if update.AuthorId != "" && len(tags) > 0 {
			update.Tags = tags

			jsonifiedMsg, err := json.Marshal(update)
			if err != nil {
				slog.Warn(err.Error())
				return
			}

			err = receiver.kafkaHelper.SendToTopic(receiver.conf.Kafka.Topic, []byte(update.Hash), jsonifiedMsg)
			if err != nil {
				slog.Warn(err.Error())
				return
			}
			slog.Info(fmt.Sprintf("Send message: %s", update.Hash))
		}
	}
}
