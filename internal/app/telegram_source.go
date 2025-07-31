package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"strconv"
	"wombat/internal/domain"
)

type TelegramSource struct {
	*tgbotapi.BotAPI
	fwdChan chan *domain.Message
}

func NewTelegramSource(token string, fwdChan chan *domain.Message) (*TelegramSource, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &TelegramSource{
		BotAPI:  bot,
		fwdChan: fwdChan,
	}, nil
}

func (receiver *TelegramSource) GetSourceType() domain.SourceType {
	return domain.TELEGRAM
}

func (receiver *TelegramSource) Process() {
	u := tgbotapi.NewUpdate(0)
	u.AllowedUpdates = append(
		u.AllowedUpdates,
		tgbotapi.UpdateTypeMessageReaction,
		tgbotapi.UpdateTypeMessage,
		tgbotapi.UpdateTypeEditedMessage,
	)
	u.Timeout = 60

	slog.Info(fmt.Sprintf("Authorized on account %s", receiver.Self.UserName))

	for update := range receiver.GetUpdatesChan(u) {

		processMessage := func(message *tgbotapi.Message) {
			if message.From.UserName == "" {
				slog.Warn("Missed username")
				return
			}

			messageId := strconv.Itoa(message.MessageID)
			chatId := strconv.FormatInt(message.Chat.ID, 10)
			msg := &domain.Message{
				TargetType: domain.JIRA,
				SourceType: domain.TELEGRAM,
				Content:    message.Text,
				UserId:     message.From.UserName,
				ChatId:     chatId,
				MessageId:  messageId,
			}
			receiver.fwdChan <- msg
		}

		if update.Message != nil {
			processMessage(update.Message)
		}
		if update.EditedMessage != nil {
			processMessage(update.EditedMessage)
		}
	}
}
