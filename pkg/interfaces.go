package pkg

import "context"

type Plugin interface {
	Init([]byte) error
	IsInit() bool
}

// Producer describes the producer type.
type Producer interface {
	Close() error

	// Publish provides channel for reading source incoming data.
	Publish(ctx context.Context) <-chan []byte
}

// Consumer describes the consumer type.
type Consumer interface {
	Close() error

	// Consume provided.
	// args is json representation of the expected structure.
	Consume(args []byte) error
}

// Transform provides filtration/transformation mechanism.
type Transform interface {

	/*
	   Eval evaluates provided object either to the next types:
	       - basic Go types (bool, string, int, etc.);
	       - object in JSON format.
	*/
	Eval(obj []byte) (any, error)
}
