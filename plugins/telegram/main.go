package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/alersrt/wombat/pkg"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	Token string `yaml:"token"`
}

type Plugin struct {
	mtx    sync.Mutex
	isInit atomic.Bool
	cfg    *Config
	bot    *api.BotAPI
}

func Export() pkg.Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(cfg []byte) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.cfg = &Config{}
	if err := json.Unmarshal(cfg, p.cfg); err != nil {
		return err
	}

	bot, err := api.NewBotAPI(p.cfg.Token)
	if err != nil {
		return fmt.Errorf("tg: new: %v", err)
	}
	p.bot = bot

	p.isInit.Store(true)
	return nil
}

func (p *Plugin) IsInit() bool {
	return p.isInit.Load()
}

func (p *Plugin) Close() error {
	return nil
}

type SendArgs struct {
	ChatId  int64  `json:"chat_id"`
	Content string `json:"content"`
}

func (p *Plugin) Consume(args []byte) error {
	if !p.IsInit() {
		return fmt.Errorf("tg: send: not init")
	}

	sA := &SendArgs{}
	if err := json.Unmarshal(args, sA); err != nil {
		return err
	}
	_, err := p.bot.Send(api.NewMessage(sA.ChatId, sA.Content))
	return err
}
