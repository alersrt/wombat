package source

type MockSource struct {
	Source chan any
	Done   chan bool
}

func (receiver *MockSource) ForwardTo(target chan any) {
	receiver.Done <- true
	for update := range receiver.Source {
		target <- update
	}
}
