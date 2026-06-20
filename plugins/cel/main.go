package cel

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alersrt/wombat/pkg"
)

type Config struct {
	Expr string `yaml:"expr"`
}

type Plugin struct {
	mtx    sync.Mutex
	isInit atomic.Bool
	cfg    *Config
	expr   *Cel
}

func Export() pkg.Component {
	return &Plugin{}
}

func (p *Plugin) Type() pkg.ComponentType {
	return pkg.ComponentTypeProcessor
}

func (p *Plugin) Init(cfg pkg.Config) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.cfg = &Config{}
	if err := json.Unmarshal(cfg.Value, p.cfg); err != nil {
		return err
	}

	expr, err := NewCel(p.cfg.Expr)
	if err != nil {
		return fmt.Errorf("cel: new: %v", err)
	}
	p.expr = expr

	p.isInit.Store(true)
	return nil
}

func (p *Plugin) IsInit() bool {
	return p.isInit.Load()
}

func (p *Plugin) Close() error {
	p.isInit.Store(false)
	return nil
}

func (p *Plugin) Process(ctx context.Context, input <-chan pkg.Message) (<-chan pkg.Message, error) {
	if !p.IsInit() {
		return nil, fmt.Errorf("cel: run: not init")
	}

	response := make(chan pkg.Message)
	defer close(response)

	go func() {
		for p.IsInit() {
			select {
			case <-ctx.Done():
				return
			case req := <-input:
				res, err := p.expr.EvalBytes(req.Value)
				if err != nil {
					slog.Warn(fmt.Sprintf("cel: run: %v", err))
					continue
				}

				response <- pkg.Message{
					Headers:   req.Headers,
					Value:     res,
					Timestamp: time.Now(),
				}
			}
		}
	}()

	return response, nil
}
