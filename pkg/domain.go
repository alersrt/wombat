package pkg

import (
	"time"
)

type ComponentType string

const (
	ComponentTypeSource    ComponentType = "source"
	ComponentTypeProcessor ComponentType = "processor"
	ComponentTypeSink      ComponentType = "sink"
	ComponentTypeLookup    ComponentType = "lookup"
)

type Message struct {
	Headers   map[string][]byte `json:"headers"`
	Key       []byte            `json:"key"`
	Value     []byte            `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Error     error             `json:"error"`
}

type Config struct {
	Value []byte `json:"value"`
}
