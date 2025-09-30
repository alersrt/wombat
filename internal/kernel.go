package internal

import (
	"encoding/json"
	"fmt"
	"plugin"

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

	producers := make(map[string]pkg.Producer)
	for _, p := range cfg.Producers {
		cfg, err := json.Marshal(p.Conf)
		if err != nil {
			return nil, fmt.Errorf("kernel: producer: [%s]: %w", p.Name, err)
		}

		if plugins[p.Plugin] == nil {
			return nil, fmt.Errorf("kernel: producer: [%s]: not exist", p.Name)
		}
		dumb := plugins[p.Plugin]()
		err = dumb.Init(cfg)
		if err != nil {
			return nil, fmt.Errorf("kernel: producer: [%s]: %w", p.Name, err)
		}
		producers[p.Name] = dumb.(pkg.Producer)
	}

	consumers := make(map[string]pkg.Consumer)
	for _, c := range cfg.Consumers {
		cfg, err := json.Marshal(c.Conf)
		if err != nil {
			return nil, fmt.Errorf("kernel: consumer: [%s]: %w", c.Name, err)
		}

		if plugins[c.Plugin] == nil {
			return nil, fmt.Errorf("kernel: producer: [%s]: not exist", c.Name)
		}
		dumb := plugins[c.Plugin]()
		err = dumb.Init(cfg)
		if err != nil {
			return nil, fmt.Errorf("kernel: consumer: [%s]: %w", c.Name, err)
		}
		consumers[c.Name] = dumb.(pkg.Consumer)
	}

	fmt.Printf("!!!!!!!!!!!!!!!!!!!!!!")

	return &Kernel{}, nil
}
