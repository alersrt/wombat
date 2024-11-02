package config

type MockConfig struct {
	*Config
}

func (receiver *MockConfig) Init(args []string) error {
	receiver.isInitiated = true
	return nil
}

func (receiver *MockConfig) IsInitiated() bool {
	return true
}
