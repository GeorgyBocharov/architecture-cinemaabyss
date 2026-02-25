package app

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func registerKafkaProducer(config *kafka.ConfigMap) (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	go func() {
		for e:= range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message: 
				if ev.TopicPartition.Error != nil {
					topicPartition := ""
					if ev.TopicPartition.Topic != nil {
						topicPartition = fmt.Sprintf("%s:%d-%d", *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
					}
					log.Printf("Error delivery to topic %s: %v", topicPartition, ev.TopicPartition.Error)
				}
			case kafka.Error: 
				log.Printf("Kafka error in producer: %v", ev)
			}
		}
	}()

	return producer, err
}