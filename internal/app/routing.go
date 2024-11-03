package app

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"log/slog"
	"regexp"
	"strconv"
	"wombat/internal/domain"
)

func (receiver *Application) route() {

	hash := func(sourceType domain.SourceType, chatId string, messageId string) string {
		return uuid.NewSHA1(uuid.NameSpaceURL, []byte(sourceType.String()+chatId+messageId)).String()
	}

	for update := range receiver.routeChan {
		switch matched := update.(type) {
		case tgbotapi.Update:
			if matched.Message != nil {
				pattern := regexp.MustCompile(receiver.conf.Bot.Tag)
				tags := pattern.FindAllString(matched.Message.Text, -1)

				if len(tags) > 0 {
					messageId := strconv.Itoa(matched.Message.MessageID)
					chatId := strconv.FormatInt(matched.Message.Chat.ID, 10)

					msg := &domain.MessageEvent{
						Hash:       hash(domain.TELEGRAM, chatId, messageId),
						EventType:  domain.CREATE,
						SourceType: domain.TELEGRAM,
						Text:       matched.Message.Text,
						AuthorId:   matched.Message.From.UserName,
						ChatId:     chatId,
						MessageId:  messageId,
					}

					jsonifiedMsg, err := json.Marshal(msg)
					if err != nil {
						slog.Warn(err.Error())
						return
					}

					err = receiver.kafkaHelper.SendToTopic(receiver.conf.Kafka.Topic, jsonifiedMsg)
					if err != nil {
						slog.Warn(err.Error())
						return
					}
					slog.Info(fmt.Sprintf("Sent message: %s", tags))
				}
			}
		}
	}
}
