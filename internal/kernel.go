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

func NewKernel(plugs []*PluginCfg) (*Kernel, error) {
	srcMap := make(map[string]pkg.Producer)
	dstMap := make(map[string]pkg.Consumer)
	celMap := make(map[string]pkg.Transform)
	for _, p := range plugs {
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
			srcMap[p.Name] = pl
		case pkg.Consumer:
			dstMap[p.Name] = pl
		case pkg.Transform:
			celMap[p.Name] = pl
		default:
			return nil, fmt.Errorf("kernel: unknown interface [%s]", p.Name)
		}
	}

	return &Kernel{producers: srcMap, consumers: dstMap, cel: celMap}, nil
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
