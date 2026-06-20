package internal

import (
	"context"
	"fmt"
	"log/slog"
	"plugin"
	"wombat/internal/cel"

	"github.com/alersrt/wombat/pkg"
)

type Kernel struct {
	broadcasts map[string]BroadcastServer
	routes     map[string]*route
}

type route struct {
	from pkg.Component
	to   pkg.Component
	when *cel.Cel
	expr *cel.Cel
}

func (r *route) Close() error {
	r.from.Close()
	r.to.Close()
	return nil
}

func (k *Kernel) Serve(ctx context.Context) {

	<-ctx.Done()
}

func NewKernel(cfg *Config) (*Kernel, error) {
	plugins := make(map[string]func() pkg.Component)
	for _, p := range cfg.Components {
		open, err := plugin.Open(p.Plugin.Path)
		if err != nil {
			return nil, fmt.Errorf("kernel: plugin: [%s]: %w", p.ID, err)
		}
		lookup, err := open.Lookup("Export")
		if err != nil {
			return nil, fmt.Errorf("kernel: plugin: [%s]: %w", p.ID, err)
		}

		plugins[p.ID] = lookup.(func() pkg.Component)
	}

	slog.Info("kernel: ok")

	return &Kernel{broadcasts: nil, routes: nil}, nil
}
