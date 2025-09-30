package internal

import (
	"encoding/json"
	"fmt"
	"plugin"
	"wombat/internal/broadcast"

	"github.com/alersrt/wombat/pkg"
)

type Kernel struct {
}

type Executor struct {
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

	broadcasts := make(map[string]broadcast.BroadcastServer)
	producers := make(map[string]pkg.Producer)
	for _, p := range cfg.Producers {
		cfg, err := json.Marshal(p.Conf)
		if err != nil {
			return nil, fmt.Errorf("kernel: producer: [%s]: %w", p.Name, err)
		}

		if plugins[p.Plugin] == nil {
			return nil, fmt.Errorf("kernel: producer: plugin: [%s]: not exist", p.Plugin)
		}
		dumb := plugins[p.Plugin]()
		err = dumb.Init(cfg)
		if err != nil {
			return nil, fmt.Errorf("kernel: producer: [%s]: %w", p.Name, err)
		}
		producer, ok := dumb.(pkg.Producer)
		if !ok {
			return nil, fmt.Errorf("kernel: producer: [%s]: not implemented", p.Name)
		}
		producers[p.Name] = producer
		broadcasts[p.Name] = broadcast.NewBroadcastServer(producer.Publish())
	}

	consumers := make(map[string]pkg.Consumer)
	for _, c := range cfg.Consumers {
		cfg, err := json.Marshal(c.Conf)
		if err != nil {
			return nil, fmt.Errorf("kernel: consumer: [%s]: %w", c.Name, err)
		}

		if plugins[c.Plugin] == nil {
			return nil, fmt.Errorf("kernel: consumer: plugin: [%s]: not exist", c.Plugin)
		}
		dumb := plugins[c.Plugin]()
		err = dumb.Init(cfg)
		if err != nil {
			return nil, fmt.Errorf("kernel: consumer: [%s]: %w", c.Name, err)
		}
		consumer, ok := dumb.(pkg.Consumer)
		if !ok {
			return nil, fmt.Errorf("kernel: consumer: [%s]: not implemented", c.Name)
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

	fmt.Printf("!!!!!!!!!!!!!!!!!!!!!!")

	return &Kernel{}, nil
}
