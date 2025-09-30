package internal

import (
	"fmt"
	"plugin"
	"github.com/alersrt/wombat/pkg"
)

type Kernel struct {
}

type Executor struct {
}

func NewKernel(cfg *Config) (*Kernel, error) {
	plugins := make(map[string]pkg.Plugin)
	for _, p := range cfg.Plugins {
		open, err := plugin.Open(p.Bin)
		if err != nil {
			return nil, fmt.Errorf("kernel: plugin: [%s]: %v", p.Name, err)
		}
		lookup, err := open.Lookup("Export")
		if err != nil {
			return nil, fmt.Errorf("kernel: plugin: [%s]: %v", p.Name, err)
		}

		plugins[p.Name] = lookup.(func() pkg.Plugin)()
	}

    fmt.Printf("!!!!!!!!!!!!!!!!!!!!!!")

	return &Kernel{}, nil
}
