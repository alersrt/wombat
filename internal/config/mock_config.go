package config

type MockConfig Config

func (receiver *MockConfig) Init(args []string) error {
	receiver.Bot.Tag = "(TEST-\\d+)"
	receiver.Bot.Emoji = "👍"
	receiver.Kafka.Topic = "wombat.test.topic"
	receiver.Kafka.Bootstrap = "localhost:9092"
	receiver.Kafka.GroupId = "wombat-test"
	receiver.isInitiated = true
	return nil
}

func (receiver *MockConfig) IsInitiated() bool {
	return true
}
