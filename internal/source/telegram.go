package source

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"wombat/internal/config"
	"wombat/pkg/errors"
)

type telegram struct {
	Updates chan any
	Config  *config.Config
}

func NewTelegramSource(updates chan any, config *config.Config) Source {
	return &telegram{
		Updates: updates,
		Config:  config,
	}
}

func (receiver *telegram) Read(cancel context.CancelCauseFunc) {
	if receiver.Updates == nil {
		cancel(errors.NewError("Updates channel is not defined"))
	}
	bot, err := tgbotapi.NewBotAPI(receiver.Config.Telegram.Token)
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
