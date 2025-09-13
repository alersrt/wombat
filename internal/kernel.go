package internal

import (
	"encoding/json"
	"plugin"
	"wombat/pkg"
)

type Kernel struct {
	src map[string]pkg.Src
	dst map[string]pkg.Dst
}

func NewKernel(src []*PluginCfg, dst []*PluginCfg) (*Kernel, error) {
	srcMap := make(map[string]pkg.Src)
	dstMap := make(map[string]pkg.Dst)
	for _, p := range src {
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
			nS, err := newPlug(bytes)
			if err != nil {
				return nil, err
			}
			srcMap[p.Name] = nS
		}

	}

	return &Kernel{src: srcMap, dst: dstMap}, nil
}
