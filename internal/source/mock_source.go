package source

type MockSource struct {
	Source chan any
}

func (receiver *MockSource) ForwardTo(target chan any) {
	for update := range receiver.Source {
		target <- update
	}
}
