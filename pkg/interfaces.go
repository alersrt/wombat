package pkg

import "context"

type Plugin interface {
	Init([]byte) error
	IsInit() bool
}

// Src describes source type.
type Src interface {
	Close() error

	// Run runs source listening until ctx done or error happened.
	Run(ctx context.Context) error

	// Publish provides channel for reading source incoming data.
	Publish() <-chan []byte
}

// Dst describes destination type.
type Dst interface {
	Close() error

	// Send data to destination.
	// args is json representation of the expected structure.
	Send(args []byte) error
}

// Cel provides filtration/transformation mechanism.
type Cel interface {

	/*
	   Eval evaluates provided object either to the next types:
	       - basic Go types (bool, string, int, etc.);
	       - object in JSON format.
	*/
	Eval(obj []byte) (any, error)
}
