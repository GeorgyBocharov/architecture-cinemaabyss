package app

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Container struct {
	PaymentsConsumer *kafka.Consumer
	MoviesConsumer *kafka.Consumer
	UserConsumer *kafka.Consumer

	PaymentsProducer *kafka.Producer
	MoviesProducer *kafka.Producer
	UserProducer *kafka.Producer
}