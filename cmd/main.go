package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"wombat/internal/config"
	"wombat/pkg/daemon"
	log "wombat/pkg/log"
)

func main() {
	conf := &config.Config{}
	daemon.Create(conf, func(cancel context.CancelCauseFunc) {
		bot, err := tgbotapi.NewBotAPI(conf.Telegram.Token)
		if err != nil {
			log.Error(err)
			cancel(err)
		}

		bot.Debug = true

		log.Infof("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.AllowedUpdates = append(u.AllowedUpdates, tgbotapi.UpdateTypeMessageReaction, tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeEditedMessage)
		u.Timeout = 60

		updates := bot.GetUpdatesChan(u)

		for update := range updates {

			if update.EditedMessage != nil {
				pattern := regexp.MustCompile(conf.Bot.Tag)
				tag := pattern.FindAllString(update.EditedMessage.Text, -1)
				log.Info(tag)
				log.Info(update.EditedMessage.Text)
			}
		}
	})
}
