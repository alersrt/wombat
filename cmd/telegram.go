package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type Telegram struct {
	Updates chan any
}

func (receiver *Telegram) Read(cancel context.CancelCauseFunc) {
	bot, err := tgbotapi.NewBotAPI(conf.Telegram.Token)
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
