package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"wombat/internal/config"
	"wombat/pkg/daemon"
	. "wombat/pkg/log"
)

func main() {
	conf := &config.Config{}
	daemon.Create(conf, func() {
		bot, err := tgbotapi.NewBotAPI(conf.Telegram.Token)
		if err != nil {
			ErrorLog.Fatal(err)
		}

		bot.Debug = true

		InfoLog.Printf("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.AllowedUpdates = append(u.AllowedUpdates, tgbotapi.UpdateTypeMessageReaction, tgbotapi.UpdateTypeMessage)
		u.Timeout = 60

		updates := bot.GetUpdatesChan(u)

		for update := range updates {

			if update.MessageReaction != nil {
				InfoLog.Printf("[%s] %s", update.MessageReaction.User.UserName, update.MessageReaction.NewReaction)
				msg := tgbotapi.NewCopyMessage(update.MessageReaction.Chat.ID, update.MessageReaction.Chat.ID, update.MessageReaction.MessageID)
				_, err := bot.Send(msg)
				if err != nil {
					InfoLog.Println(err)
				}
			}

			if update.Message != nil {
				InfoLog.Println(update.Message.Text)
			}
		}
	})
}
