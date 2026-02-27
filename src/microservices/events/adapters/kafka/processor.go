package kafka

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type (
	KafkaProcessor[T any] struct {
		processor Processor[T]
	}
	Processor[T any] interface {
		Process(ctx context.Context, value T) error
	}
)

func NewKafkaProcessor[T any](processor Processor[T]) *KafkaProcessor[T] {
	return &KafkaProcessor[T]{processor: processor}
}

func (p *KafkaProcessor[T]) Process(ctx context.Context, message *kafka.Message) error {
	var val T
	if err := json.Unmarshal(message.Value, &val); err != nil {
		return fmt.Errorf("failed to umarshall value as json: %v", err)
	}
	
	return p.processor.Process(ctx, val)
}