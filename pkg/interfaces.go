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
	Publish() <-chan []byte

	// Run starts reading data from source and publish it to channel.
	Run(ctx context.Context) error
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
		Eval evaluates provided object to bytes of JSON object.
		It can be not only struct but also bool, string, int, etc..
	*/
	EvalBytes(obj []byte) ([]byte, error)
}
