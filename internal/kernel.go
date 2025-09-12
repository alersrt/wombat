package internal

import "sync"

type Source interface {
}

type Target interface {
}

type Action interface {
	Execute(args ...any)
}

type Kernel struct {
	rwMtx sync.RWMutex
}
