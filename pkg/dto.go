package pkg

import "time"

type Args[T any] struct {
	Headers   map[string]string `json:"headers"`
	Value     T                 `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
}
