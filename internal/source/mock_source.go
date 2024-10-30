package source

import "context"

type MockSource struct {
	Updates chan any
	Source  chan any
	Done    chan bool
}

func (receiver *MockSource) Read(causeFunc context.CancelCauseFunc) {
	receiver.Done <- true
	for update := range receiver.Source {
		receiver.Updates <- update
	}
}
