package internal

import (
	"context"
	"fmt"
	"plugin"
)

type Source interface {
	Close() error
	Read(ctx context.Context) error
	Publish() <-chan *Message
	Name() string
}

type Target interface {
	Close() error
}

type Action interface {
	Execute(args ...any) error
}

type Kernel struct {
	actions map[string]Action
}

func NewKernel(plugs []*Plugin) (*Kernel, error) {
	actions := make(map[string]Action)
	for _, p := range plugs {
		open, err := plugin.Open(p.Path)
		if err != nil {
			return nil, err
		}

		lookup, err := open.Lookup("Action")
		if err != nil {
			return nil, err
		}

		action, ok := lookup.(Action)
		if !ok {
			return nil, fmt.Errorf("kernel: new: wrong plugin name=%s", p.Name)
		}

		actions[p.Name] = action
	}

	return &Kernel{actions: actions}, nil
}
