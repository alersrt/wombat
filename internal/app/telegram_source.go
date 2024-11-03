package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"strconv"
	"wombat/internal/domain"
)

type TelegramSource struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramSource(token string) (*TelegramSource, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &TelegramSource{bot: bot}, nil
}

func (receiver *TelegramSource) ForwardTo(target chan *domain.MessageEvent) {
	u := tgbotapi.NewUpdate(0)
	u.AllowedUpdates = append(
		u.AllowedUpdates,
		tgbotapi.UpdateTypeMessageReaction,
		tgbotapi.UpdateTypeMessage,
		tgbotapi.UpdateTypeEditedMessage,
	)
	u.Timeout = 60

	slog.Info(fmt.Sprintf("Authorized on account %s", receiver.bot.Self.UserName))

	for update := range receiver.bot.GetUpdatesChan(u) {

		if update.Message != nil {

			messageId := strconv.Itoa(update.Message.MessageID)
			chatId := strconv.FormatInt(update.Message.Chat.ID, 10)

			msg := &domain.MessageEvent{
				Hash:       domain.Hash(domain.TELEGRAM, chatId, messageId),
				SourceType: domain.TELEGRAM,
				Text:       update.Message.Text,
				AuthorId:   update.Message.From.UserName,
				ChatId:     chatId,
				MessageId:  messageId,
			}

			target <- msg
		}
	}
}
