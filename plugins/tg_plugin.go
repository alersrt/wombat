package main

import (
	"encoding/json"
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"wombat/internal"
)

type TgCfg struct {
	Token string `yaml:"token"`
}

type TgDst struct {
	cfg *TgCfg
	bot *api.BotAPI
}

func New(cfg map[string]any) (*TgDst, error) {
	bytes, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	tgCfg := &TgCfg{}
	if err := json.Unmarshal(bytes, tgCfg); err != nil {
		return nil, err
	}

	bot, err := api.NewBotAPI(tgCfg.Token)
	if err != nil {
		return nil, fmt.Errorf("tg: new: %v", err)
	}
	return &TgDst{cfg: tgCfg, bot: bot}, nil
}

func (t *TgDst) Close() error {
	return nil
}

func (t *TgDst) Send(chatId int64, msg string) error {
	_, err := t.bot.Send(api.NewMessage(chatId, msg))
	return err
}
