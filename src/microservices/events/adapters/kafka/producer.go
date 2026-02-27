package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer[T any] struct {
	topic string
	producer *kafka.Producer
	keyExtractor func(value T) []byte
}

func NewProducer[T any](topic string, producer *kafka.Producer, keyExtractor func(value T) []byte) *Producer[T] {
	return &Producer[T]{
		topic: topic,
		producer: producer,
		keyExtractor: keyExtractor,
	}
}

func (p *Producer[T]) Send(value T) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to transform value to json: %v", err)
	}
	key := p.keyExtractor(value)

	log.Printf("sending message to topic %s with key %s, value %s", p.topic, string(key), string(valueBytes))

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &p.topic,
			Partition: kafka.PartitionAny,
		},
		Key: key,
		Value: valueBytes,
	}
	err = p.producer.Produce(message, nil)
	if err != nil {
		fmt.Printf("failed to send: %v", err)
	}
	remaining := p.producer.Flush(5000)
	fmt.Printf("remaining after flush: %v\n", remaining)

	return err
}