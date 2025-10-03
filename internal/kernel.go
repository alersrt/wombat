package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"plugin"

	"github.com/alersrt/wombat/pkg"
)

type Kernel struct {
	broadcasts map[string]BroadcastServer
	rules      map[string]*rule
}

type rule struct {
	proc          *Processor
	broadcastName string
}

func (r *rule) Close() error {
	return r.proc.Close()
}

func (k *Kernel) Serve(ctx context.Context) {
	for n, v := range k.broadcasts {
		slog.Info(fmt.Sprintf("serve: broadcast: [%s]: ok", n))
		go v.Serve(ctx)
		slog.Info(fmt.Sprintf("serve: broadcast: [%s]: ok", n))
	}
	for n, v := range k.rules {
		go func() {
			defer func() {
				_ = v.Close()
			}()
			subscription := k.broadcasts[v.broadcastName].Subscribe()
			defer k.broadcasts[v.broadcastName].CancelSubscription(subscription)

			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-subscription:
					if !ok {
						slog.Warn(fmt.Sprintf("serve: rule: [%s]: stopped", n))
						return
					}
					if err := v.proc.Process(val); err != nil {
						slog.Warn(fmt.Sprintf("serve: rule: [%s]: consume: %+v", n, err))
					}
				}
			}
		}()
		slog.Info(fmt.Sprintf("serve: rule: [%s]: ok", n))
	}

	<-ctx.Done()
}

func NewKernel(cfg *Config) (*Kernel, error) {
	plugins := make(map[string]func() pkg.Plugin)
	for _, p := range cfg.Plugins {
		open, err := plugin.Open(p.Bin)
		if err != nil {
			return nil, fmt.Errorf("kernel: plugin: [%s]: %w", p.Name, err)
		}
		lookup, err := open.Lookup("Export")
		if err != nil {
			return nil, fmt.Errorf("kernel: plugin: [%s]: %w", p.Name, err)
		}

		plugins[p.Name] = lookup.(func() pkg.Plugin)
	}

	producers := make(map[string]pkg.Producer)
	broadcasts := make(map[string]BroadcastServer)
	for _, p := range cfg.Producers {
		producer, err := producer(p, plugins)
		if err != nil {
			return nil, err
		}
		producers[p.Name] = producer
		broadcasts[p.Name] = NewBroadcastServer(producer)
	}

	consumers := make(map[string]pkg.Consumer)
	for _, c := range cfg.Consumers {
		consumer, err := consumer(c, plugins)
		if err != nil {
			return nil, err
		}
		consumers[c.Name] = consumer
	}

	applies := make(map[string]*rule)
	for _, r := range cfg.Rules {
		if producers[r.Producer] == nil {
			return nil, fmt.Errorf("kernel: producer: [%s]: not exist", r.Producer)
		}
		if consumers[r.Consumer] == nil {
			return nil, fmt.Errorf("kernel: consumer: [%s]: not exist", r.Consumer)
		}
		proc, err := NewProcessor(consumers[r.Consumer], r.Filter, r.Transform)
		if err != nil {
			return nil, fmt.Errorf("kernel: rule: [%s]: %w", r.Name, err)
		}
		applies[r.Name] = &rule{
			broadcastName: r.Producer,
			proc:          proc,
		}
	}

	slog.Info("kernel: ok")

	return &Kernel{broadcasts: broadcasts, rules: applies}, nil
}

func producer(value *ItemCfg, plugins map[string]func() pkg.Plugin) (pkg.Producer, error) {
	cfg, err := json.Marshal(value.Conf)
	if err != nil {
		return nil, fmt.Errorf("kernel: producer: [%s]: %w", value.Name, err)
	}

	if plugins[value.Plugin] == nil {
		return nil, fmt.Errorf("kernel: producer: plugin: [%s]: not exist", value.Plugin)
	}
	dumb := plugins[value.Plugin]()
	err = dumb.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("kernel: producer: [%s]: %w", value.Name, err)
	}
	producer, ok := dumb.(pkg.Producer)
	if !ok {
		return nil, fmt.Errorf("kernel: producer: [%s]: not implemented", value.Name)
	}
	return producer, nil
}

func consumer(value *ItemCfg, plugins map[string]func() pkg.Plugin) (pkg.Consumer, error) {
	cfg, err := json.Marshal(value.Conf)
	if err != nil {
		return nil, fmt.Errorf("kernel: consumer: [%s]: %w", value.Name, err)
	}

	if plugins[value.Plugin] == nil {
		return nil, fmt.Errorf("kernel: consumer: plugin: [%s]: not exist", value.Plugin)
	}
	dumb := plugins[value.Plugin]()
	err = dumb.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("kernel: consumer: [%s]: %w", value.Name, err)
	}
	consumer, ok := dumb.(pkg.Consumer)
	if !ok {
		return nil, fmt.Errorf("kernel: consumer: [%s]: not implemented", value.Name)
	}
	return consumer, nil
}
