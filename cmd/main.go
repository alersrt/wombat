package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"wombat/internal/config"
	"wombat/pkg/daemon"
	log "wombat/pkg/log"
)

func main() {
	conf := &config.Config{}
	daemon.Create(conf, func(cancel context.CancelCauseFunc) {
		bot, err := tgbotapi.NewBotAPI(conf.Telegram.Token)
		if err != nil {
			log.Errorln(err)
			cancel(err)
		}

		bot.Debug = true

		log.Infof("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.AllowedUpdates = append(u.AllowedUpdates, tgbotapi.UpdateTypeMessageReaction, tgbotapi.UpdateTypeMessage)
		u.Timeout = 60

		updates := bot.GetUpdatesChan(u)

		for update := range updates {

			if update.MessageReaction != nil {
				log.Infof("[%s] %s", update.MessageReaction.User.UserName, update.MessageReaction.NewReaction)
				msg := tgbotapi.NewCopyMessage(update.MessageReaction.Chat.ID, update.MessageReaction.Chat.ID, update.MessageReaction.MessageID)
				_, err := bot.Send(msg)
				if err != nil {
					log.Warn(err)
				}
			}

			if update.Message != nil {
				log.Infoln(update.Message.Text)
			}
		}
	})
}
