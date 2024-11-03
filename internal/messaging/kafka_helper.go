package messaging

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"log/slog"
)

type MessageHelper interface {
	SendToTopic(topic string, message []byte) error
	Subscribe(topics []string, handler func(*kafka.Message) error) error
}

func NewKafkaHelper(configMap *kafka.ConfigMap) (MessageHelper, error) {
	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &KafkaMessageHelper{
		producer: producer,
		consumer: consumer,
	}, nil
}

type KafkaMessageHelper struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
}

func (receiver *KafkaMessageHelper) SendToTopic(topic string, message []byte) error {
	key, err := uuid.New().MarshalText()
	if err != nil {
		slog.Warn(err.Error())
		return err
	}
	return receiver.producer.Produce(&kafka.Message{
		Key:   key,
		Value: message,
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
	}, nil)
}

func (receiver *KafkaMessageHelper) Subscribe(topics []string, handler func(*kafka.Message) error) error {
	err := receiver.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer receiver.consumer.Close()

	for {
		msg, err := receiver.consumer.ReadMessage(-1)
		if err != nil {
			slog.Warn(err.Error())
			continue
		}
		if err := handler(msg); err != nil {
			slog.Warn(err.Error())
			continue
		}
	}
}
