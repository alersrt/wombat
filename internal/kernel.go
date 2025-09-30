package internal

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"plugin"
	"wombat/internal/broadcast"

	"github.com/alersrt/wombat/pkg"
)

type Kernel struct {
	broadcasts map[string]broadcast.BroadcastServer
	producers  map[string]pkg.Producer
}

func NewKernel(cfg *Config) (*Kernel, error) {
	plugins := make(map[string]func() pkg.Plugin)
	for _, p := range cfg.Plugins {
		open, err := plugin.Open(p.Bin)
		if err != nil {
			return nil, fmt.Errorf("kernel: plugin: [%s]: %w", p.Name, err)
		}
		lookup, err := open.Lookup("Export")
		if err != nil {
			return nil, fmt.Errorf("kernel: plugin: [%s]: %w", p.Name, err)
		}

		plugins[p.Name] = lookup.(func() pkg.Plugin)
	}

	producers := make(map[string]pkg.Producer)
	broadcasts := make(map[string]broadcast.BroadcastServer)
	for _, p := range cfg.Producers {
		producer, err := producer(p, plugins)
		if err != nil {
			return nil, err
		}
		producers[p.Name] = producer
		broadcasts[p.Name] = broadcast.NewBroadcastServer(producer.Publish())
	}

	consumers := make(map[string]pkg.Consumer)
	for _, c := range cfg.Consumers {
		consumer, err := consumer(c, plugins)
		if err != nil {
			return nil, err
		}
		consumers[c.Name] = consumer
	}

	for _, r := range cfg.Rules {
		if producers[r.Producer] == nil {
			return nil, fmt.Errorf("kernel: producer: [%s]: not exist", r.Producer)
		}
		if consumers[r.Consumer] == nil {
			return nil, fmt.Errorf("kernel: consumer: [%s]: not exist", r.Consumer)
		}

	}

	slog.Info("!!!!!!!!!!!!!!!!!!!!!11")

	return &Kernel{}, nil
}

func producer(value *ItemCfg, plugins map[string]func() pkg.Plugin) (pkg.Producer, error) {
	cfg, err := json.Marshal(value.Conf)
	if err != nil {
		return nil, fmt.Errorf("kernel: producer: [%s]: %w", value.Name, err)
	}

	if plugins[value.Plugin] == nil {
		return nil, fmt.Errorf("kernel: producer: plugin: [%s]: not exist", value.Plugin)
	}
	dumb := plugins[value.Plugin]()
	err = dumb.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("kernel: producer: [%s]: %w", value.Name, err)
	}
	producer, ok := dumb.(pkg.Producer)
	if !ok {
		return nil, fmt.Errorf("kernel: producer: [%s]: not implemented", value.Name)
	}
	return producer, nil
}

func consumer(value *ItemCfg, plugins map[string]func() pkg.Plugin) (pkg.Consumer, error) {
	cfg, err := json.Marshal(value.Conf)
	if err != nil {
		return nil, fmt.Errorf("kernel: consumer: [%s]: %w", value.Name, err)
	}

	if plugins[value.Plugin] == nil {
		return nil, fmt.Errorf("kernel: consumer: plugin: [%s]: not exist", value.Plugin)
	}
	dumb := plugins[value.Plugin]()
	err = dumb.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("kernel: consumer: [%s]: %w", value.Name, err)
	}
	consumer, ok := dumb.(pkg.Consumer)
	if !ok {
		return nil, fmt.Errorf("kernel: consumer: [%s]: not implemented", value.Name)
	}
	return consumer, nil
}
