package source

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type telegram struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramSource(token string) (Source, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &telegram{bot: bot}, nil
}

func (receiver *telegram) ForwardTo(target chan any) {
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
		target <- update
	}
}
