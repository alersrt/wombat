package internal

import (
	"context"
	"fmt"
	"plugin"
)

type Src interface {
	Close() error
	Run(ctx context.Context) error
	Publish() <-chan any
}

type Dst interface {
	Close() error
	Send(...any) error
}

type Kernel struct {
	src map[string]Src
	dst map[string]Dst
}

func NewKernel(src []*PluginCfg, dst []*PluginCfg) (*Kernel, error) {
	srcMap := make(map[string]Src)
	for _, p := range src {
		open, err := plugin.Open(p.Bin)
		if err != nil {
			return nil, err
		}

		lookup, err := open.Lookup("New")
		if err != nil {
			return nil, err
		}

		newSrc, ok := lookup.(func(cfg map[string]any) (Src, error))
		if !ok {
			return nil, fmt.Errorf("kernel: new: wrong plugin [%s]", p.Name)
		}

		nS, err := newSrc(p.Conf)
		if err != nil {
			return nil, err
		}

		srcMap[p.Name] = nS
	}

	dstMap := make(map[string]Dst)
	for _, p := range dst {
		open, err := plugin.Open(p.Bin)
		if err != nil {
			return nil, err
		}

		lookup, err := open.Lookup("New")
		if err != nil {
			return nil, err
		}

		newDst, ok := lookup.(func(cfg map[string]any) (Dst, error))
		if !ok {
			return nil, fmt.Errorf("kernel: new: wrong plugin [%s]", p.Name)
		}

		nS, err := newDst(p.Conf)
		if err != nil {
			return nil, err
		}

		dstMap[p.Name] = nS
	}

	return &Kernel{src: srcMap, dst: dstMap}, nil
}
