package internal

import (
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	cfg *Telegram
	bot *api.BotAPI
}

func NewTgBot(cfg *Telegram) (*TgBot, error) {
	bot, err := api.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("tg: new: %v", err)
	}
	return &TgBot{cfg: cfg, bot: bot}, nil
}

func (t *TgBot) Close() error {
	return nil
}

func (t *TgBot) Send(chatId int64, msg string) error {
	_, err := t.bot.Send(api.NewMessage(chatId, msg))
	return err
}
