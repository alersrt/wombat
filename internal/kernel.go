package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"plugin"
	"wombat/internal/broadcast"

	"github.com/alersrt/wombat/pkg"
)

type Kernel struct {
	broadcasts map[string]broadcast.BroadcastServer
	applies    map[string]*applience
}

type applience struct {
	consumer      pkg.Consumer
	broadcastName string
}

func (s *applience) Close() error {
	return s.consumer.Close()
}

func (k *Kernel) Serve(ctx context.Context) {
	for n, v := range k.broadcasts {
		go v.Serve(ctx)
		slog.Info(fmt.Sprintf("serve: broadcast: %s", n))
	}
	for n, v := range k.applies {
		go func() {
			defer v.Close()
            subscription := k.broadcasts[v.broadcastName].Subscribe()
			for {
				select {
                case <-ctx.Done():
                    return
                case val, ok := <- subscription:
                    if !ok {
                        return;
                    }

                    if val != nil {
                        continue
                    }
                }
			}
		}()
		slog.Info(fmt.Sprintf("serve: subscribe [%s] on [%s]", n, v.broadcastName))
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
	broadcasts := make(map[string]broadcast.BroadcastServer)
	for _, p := range cfg.Producers {
		producer, err := producer(p, plugins)
		if err != nil {
			return nil, err
		}
		producers[p.Name] = producer
		broadcasts[p.Name] = broadcast.NewBroadcastServer(producer)
	}

	consumers := make(map[string]pkg.Consumer)
	for _, c := range cfg.Consumers {
		consumer, err := consumer(c, plugins)
		if err != nil {
			return nil, err
		}
		consumers[c.Name] = consumer
	}

	applies := make(map[string]*applience)
	for _, r := range cfg.Rules {
		if producers[r.Producer] == nil {
			return nil, fmt.Errorf("kernel: producer: [%s]: not exist", r.Producer)
		}
		if consumers[r.Consumer] == nil {
			return nil, fmt.Errorf("kernel: consumer: [%s]: not exist", r.Consumer)
		}
		applies[r.Name] = &applience{consumer: consumers[r.Consumer], broadcastName: r.Producer}
	}

	slog.Info("kernel: loaded")

	return &Kernel{broadcasts: broadcasts, applies: applies}, nil
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
