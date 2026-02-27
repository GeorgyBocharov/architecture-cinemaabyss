package app

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func registerKafkaConsumer(config map[string]any) (*kafka.Consumer, error) {
	configMap := toConfigMap(config)
	consumer, err := getConsumerWithRetries(10, configMap)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func registerKafkaProducer(config map[string]any) (*kafka.Producer, error) {
	configMap := toConfigMap(config)
	producer, err := getProducerWithRetries(10, configMap)
	if err != nil {
		return nil, err
	}

	go func() {
		for e:= range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message: 
			topicPartition := ""
				if ev.TopicPartition.Topic != nil {
					topicPartition = fmt.Sprintf("%s:%d-%d", *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
				if ev.TopicPartition.Error != nil {
					log.Printf("Error delivery to topic %s: %v", topicPartition, ev.TopicPartition.Error)
				} else {
					log.Printf("Success delivery to topic %s", topicPartition)
				}
			case kafka.Error: 
				log.Printf("Kafka error in producer: %v", ev)
			}
		}
	}()

	return producer, err
}

func getProducerWithRetries(maxRetries int, configMap *kafka.ConfigMap) (*kafka.Producer, error) {
    var  producer *kafka.Producer 
	var err error
	for i := 0; i < maxRetries; i++ {
      	producer, err = kafka.NewProducer(configMap)
		if err == nil {
			break
		}
		fmt.Printf("failed to create producer, attempt %d of %d: %v", i, maxRetries, err)

        time.Sleep(time.Duration(1) * time.Second)
    }
	metadata, err := producer.GetMetadata(nil, true, 30000) // 3 секунды таймаут
    if err != nil {
        producer.Close()
        return nil, fmt.Errorf("kafka недоступна: %w", err)
    }
	for name := range metadata.Topics {
		log.Printf("topic  - %s", name)
    }

	return producer, err
}

func getConsumerWithRetries(maxRetries int, configMap *kafka.ConfigMap) (*kafka.Consumer, error) {
    var  consumer *kafka.Consumer 
	var err error
	for i := 0; i < maxRetries; i++ {
      	consumer, err = kafka.NewConsumer(configMap)
		if err == nil {
			break
		}
		fmt.Printf("failed to create consumer, attempt %d of %d: %v", i, maxRetries, err)

        time.Sleep(time.Duration(1) * time.Second)
    }
	_, err = consumer.GetMetadata(nil, true, 30000) // 3 секунды таймаут
    if err != nil {
        consumer.Close()
        return nil, fmt.Errorf("kafka недоступна: %w", err)
    }

	return consumer, err
}

func toConfigMap(config map[string]any) *kafka.ConfigMap {
	configMap := &kafka.ConfigMap{}
	for k, v := range config {
		configMap.SetKey(k, v)
	}
	return configMap
}