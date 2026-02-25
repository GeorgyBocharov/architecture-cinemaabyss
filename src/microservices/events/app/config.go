package app

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Config struct {
	PaymentsTopic string
	UsersTopic    string
	MovesTopic    string

	ConsumerConfig map[string]interface{}
	ProducerConfig map[string]interface{}
}