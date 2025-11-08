package pkg

import "context"

// Processor describes pipeline with input and output.
type Processor interface {

	// Inits the processor. Takes the configurations in json representation.
	// Returns an error if some trouble is observed.
	Init([]byte) error

	// Checks if the processor initializated.
	IsInit() bool

	Close() error

	// Process consumes data, do work and return an answer.
	// Input and output are json representations of the expected structures.
	Process(ctx context.Context, input <-chan []byte) (<-chan []byte, error)
}
