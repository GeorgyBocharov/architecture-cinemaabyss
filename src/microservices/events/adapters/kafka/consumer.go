package kafka

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type (
	Consumer struct {
		topic string
		pollTimeout int
		consumer *kafka.Consumer
		processor MessageProcessor
	}
	MessageProcessor interface {
		Process(ctx context.Context, message *kafka.Message) error
	}
)

func NewConsumer(topic string, pollTimeout int, consumer *kafka.Consumer, processor MessageProcessor) *Consumer {
	return &Consumer{
		topic: topic,
		pollTimeout: pollTimeout,
		consumer: consumer,
		processor: processor,
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	defer c.consumer.Close()
	
	err := c.consumer.Subscribe(c.topic, nil) 
	if err != nil {
		return err
	}

	for {
		select {
		case <- ctx.Done():
			log.Printf("Closing consumer for topic: %s", c.topic)

			return nil
		default:
			ev := c.consumer.Poll(c.pollTimeout)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				err := c.processor.Process(ctx, e)
				if err != nil {
					log.Printf("Failed to process message from topic %s with partition: %d, offset: %d. error: %v", 
					c.topic, e.TopicPartition.Partition, e.TopicPartition.Offset, err)
				}
			case kafka.Error:
				log.Printf("Consumer error for topic %s. error: %v", c.topic, e.Error())
			default:
			}
		}
	}
}