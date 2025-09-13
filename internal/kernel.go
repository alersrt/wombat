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
	cel map[string]pkg.Cel
}

func NewKernel(plugs []*PluginCfg) (*Kernel, error) {
	srcMap := make(map[string]pkg.Src)
	dstMap := make(map[string]pkg.Dst)
	for _, p := range plugs {
		open, err := plugin.Open(p.Bin)
		if err != nil {
			return nil, err
		}

		lookup, err := open.Lookup("Export")
		if err != nil {
			return nil, err
		}

		plug, ok := lookup.(pkg.Plugin)
		if !ok {
			return nil, fmt.Errorf("kernel: incompatible plug [%s]", p.Name)
		}

		cfg, err := json.Marshal(p.Conf)
		if err != nil {
			return nil, err
		}

		if err := plug.Init(cfg); err != nil {
			return nil, err
		}
	}

	return &Kernel{src: srcMap, dst: dstMap}, nil
}
