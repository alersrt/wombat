package source

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"wombat/pkg/errors"
)

type telegram struct {
	Updates chan any
	Token   string
}

func NewTelegramSource(updates chan any, token string) Source {
	return &telegram{
		Updates: updates,
		Token:   token,
	}
}

func (receiver *telegram) Read(cancel context.CancelCauseFunc) {
	if receiver.Updates == nil {
		cancel(errors.NewError("Updates channel is not defined"))
	}
	bot, err := tgbotapi.NewBotAPI(receiver.Token)
	if err != nil {
		slog.Error(err.Error())
		cancel(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.AllowedUpdates = append(
		u.AllowedUpdates,
		tgbotapi.UpdateTypeMessageReaction,
		tgbotapi.UpdateTypeMessage,
		tgbotapi.UpdateTypeEditedMessage,
	)
	u.Timeout = 60

	slog.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	for update := range bot.GetUpdatesChan(u) {
		receiver.Updates <- update
	}
}
