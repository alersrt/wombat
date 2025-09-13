package telegram

import (
	"encoding/json"
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wombat/pkg"
)

type Config struct {
	Token string `yaml:"token"`
}

type Plugin struct {
	cfg *Config
	bot *api.BotAPI
}

func New(cfg []byte) (pkg.Dst, error) {
	tgCfg := &Config{}
	if err := json.Unmarshal(cfg, tgCfg); err != nil {
		return nil, err
	}

	bot, err := api.NewBotAPI(tgCfg.Token)
	if err != nil {
		return nil, fmt.Errorf("tg: new: %v", err)
	}
	return &Plugin{cfg: tgCfg, bot: bot}, nil
}

func (p *Plugin) Close() error {
	return nil
}

type SendArgs struct {
	ChatId int64
	Text   string
}

func (p *Plugin) Send(args []byte) error {
	sA := &SendArgs{}
	if err := json.Unmarshal(args, sA); err != nil {
		return err
	}
	_, err := p.bot.Send(api.NewMessage(sA.ChatId, sA.Text))
	return err
}
