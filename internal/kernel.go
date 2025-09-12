package internal

import (
	"context"
	"sync"
)

type Source interface {
	Close() error
	Read(ctx context.Context) error
	Publish() <-chan *Message
}

type Target interface {
	Close() error
}

type Action interface {
	Execute(args ...any)
}

type Kernel struct {
	rwMtx sync.RWMutex
}
