package internal

import (
	"encoding/json"
	"fmt"
	"plugin"
	"wombat/pkg"
)

type Kernel struct {
	src map[string]pkg.Src
	dst map[string]pkg.Dst
}

func NewKernel(plugs []*PluginCfg) (*Kernel, error) {
	srcMap := make(map[string]pkg.Src)
	dstMap := make(map[string]pkg.Dst)
	for _, p := range plugs {
		open, err := plugin.Open(p.Bin)
		if err != nil {
			return nil, err
		}

		lookup, err := open.Lookup("New")
		if err != nil {
			return nil, err
		}

		bytes, err := json.Marshal(p.Conf)
		if err != nil {
			return nil, err
		}
		switch newPlug := lookup.(type) {
		case func(cfg []byte) (pkg.Plugin, error):
			plug, err := newPlug(bytes)
			if err != nil {
				return nil, err
			}
			switch pl := plug.(type) {
			case pkg.Src:
				srcMap[p.Name] = pl
			case pkg.Dst:
				dstMap[p.Name] = pl
			default:
				return nil, fmt.Errorf("kernel: incompatible plug [%s]", p.Name)
			}
		default:
			return nil, fmt.Errorf("kernel: incompatible plug [%s]", p.Name)
		}
	}

	return &Kernel{src: srcMap, dst: dstMap}, nil
}
