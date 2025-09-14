package internal

import (
	"encoding/json"
	"fmt"
	"plugin"
	"wombat/pkg"
)

type Kernel struct {
	src map[string]pkg.Producer
	dst map[string]pkg.Consumer
	cel map[string]pkg.Transform
}

func NewKernel(plugs []*PluginCfg) (*Kernel, error) {
	srcMap := make(map[string]pkg.Producer)
	dstMap := make(map[string]pkg.Consumer)
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
		}
	}

	fmt.Printf("\nplugins:\nSRC:%v\nDST:%v\n", srcMap, dstMap)

	return &Kernel{src: srcMap, dst: dstMap}, nil
}
