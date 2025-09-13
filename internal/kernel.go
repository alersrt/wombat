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
		case func(cfg []byte) (pkg.Src, error):
			plug, err := newPlug(bytes)
			if err != nil {
				return nil, err
			}
			srcMap[p.Name] = plug
		case func(cfg []byte) (pkg.Dst, error):
			plug, err := newPlug(bytes)
			if err != nil {
				return nil, err
			}
			dstMap[p.Name] = plug
		default:
			return nil, fmt.Errorf("kernel: incompatible plug [%s]", p.Name)
		}
	}

	return &Kernel{src: srcMap, dst: dstMap}, nil
}
