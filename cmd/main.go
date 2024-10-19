package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"wombat/internal/config"
	"wombat/pkg/daemon"
)

func main() {
	conf := &config.Config{}
	daemon.Create(conf, func() {
		bot, err := tgbotapi.NewBotAPI(conf.Telegram.Token)
		if err != nil {
			log.Panic(err)
		}

		bot.Debug = true

		log.Printf("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.AllowedUpdates = append(u.AllowedUpdates, tgbotapi.UpdateTypeMessageReaction)
		u.Timeout = 60

		updates := bot.GetUpdatesChan(u)

		for update := range updates {

			if update.MessageReaction != nil {
				log.Printf("[%s] %s", update.MessageReaction.User.UserName, update.MessageReaction.NewReaction)
				msg := tgbotapi.NewMessage(update.MessageReaction.Chat.ID, "Emoji")
				msg.ReplyParameters.MessageID = update.MessageReaction.MessageID
				_, err := bot.Send(msg)
				if err != nil {
					log.Panic(err)
				}
				log.Println(update.MessageReaction)
			}
		}
	})
}
