package pkg

import "context"

// Component describes pipeline with input and output.
type Component interface {

	// Returns type of component
	Type() ComponentType

	// Takes the configurations and inits component.
	// Returns an error if some trouble is observed.
	Init(Config) error

	// Checks if the component initializated.
	IsInit() bool

	Close() error

	// Process consumes data, do work and return an answer.
	// Input and output are json representations of the expected structures.
	Process(ctx context.Context, input <-chan Message) (<-chan Message, error)
}
