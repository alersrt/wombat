package internal

import (
	"encoding/json"
	"fmt"
	"plugin"
	"wombat/pkg"
)

type Kernel struct {
	producers map[string]pkg.Producer
	consumers map[string]pkg.Consumer
	cel map[string]pkg.Transform
}

func NewKernel(cfg *Config) (*Kernel, error) {
	producersMap := make(map[string]pkg.Producer)
	consumersMap := make(map[string]pkg.Consumer)
	celMap := make(map[string]pkg.Transform)
	for _, p := range cfg.Plugins {
		open, err := plugin.Open(p.Bin)
		if err != nil {
			return nil, fmt.Errorf("kernel: [%s]: %v", p.Name, err)
		}

		lookup, err := open.Lookup("Export")
		if err != nil {
			return nil, fmt.Errorf("kernel: [%s]: %v", p.Name, err)
		}

		plug, ok := lookup.(pkg.Plugin)
		if !ok {
			return nil, fmt.Errorf("kernel: incompatible plug [%s]", p.Name)
		}

		cfg, err := json.Marshal(p.Conf)
		if err != nil {
			return nil, fmt.Errorf("kernel: [%s]: %v", p.Name, err)
		}

		if err := plug.Init(cfg); err != nil {
			return nil, fmt.Errorf("kernel: [%s]: %v", p.Name, err)
		}

		switch pl := plug.(type) {
		case pkg.Producer:
			producersMap[p.Name] = pl
		case pkg.Consumer:
			consumersMap[p.Name] = pl
		case pkg.Transform:
			celMap[p.Name] = pl
		default:
			return nil, fmt.Errorf("kernel: unknown interface [%s]", p.Name)
		}
	}

	return &Kernel{producers: producersMap, consumers: consumersMap, cel: celMap}, nil
}

func (k *Kernel) Producer(name string) pkg.Producer {
	return k.producers[name]
}

func (k *Kernel) Consumer(name string) pkg.Consumer {
	return k.consumers[name]
}

func (k *Kernel) Transform(name string) pkg.Transform {
	return k.cel[name]
}
