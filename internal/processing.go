package internal

import (
	"fmt"
	"wombat/internal/cel"
	"wombat/pkg"
)

type Processor struct {
	consumer  pkg.Consumer
	filter    *cel.Cel
	transform *cel.Cel
}

func (p *Processor) Close() error {
	return p.consumer.Close()
}

func NewProcessor(cons pkg.Consumer, filter string, transform string) (*Processor, error) {
	f, err := cel.NewCel(filter)
	if err != nil {
		return nil, fmt.Errorf("processor: %w", err)
	}
	t, err := cel.NewCel(transform)
	if err != nil {
		return nil, fmt.Errorf("processor: %w", err)
	}
	return &Processor{consumer: cons, filter: f, transform: t}, nil
}

func (p *Processor) Process(value []byte) error {
	var err error
	if p.filter != nil && !p.filter.EvalBool(value) {
		return fmt.Errorf("processor: not processed")
	}
	if p.transform != nil {
		if value, err = p.transform.EvalBytes(value); err != nil {
			return fmt.Errorf("processor: %w", err)
		}
	}

	if err = p.consumer.Consume(value); err != nil {
		return fmt.Errorf("processor: %w", err)
	}
	return nil
}
