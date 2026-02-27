package app

import (
	kafkaAdapter "events/adapters/kafka"
	"events/internal/entities"
	"events/internal/services"
	"events/server"
	"strconv"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Container struct {
	Config *Config

	PaymentsConsumer *kafkaAdapter.Consumer
	MoviesConsumer *kafkaAdapter.Consumer
	UsersConsumer *kafkaAdapter.Consumer

	Producer *kafka.Producer
	PaymentsProducer *kafkaAdapter.Producer[entities.Payment]
	MoviesProducer *kafkaAdapter.Producer[entities.Movie]
	UserProducer *kafkaAdapter.Producer[entities.User]

	PaymentsHandler HttpHandler
	MoviesHandler HttpHandler
	UsersHandler HttpHandler

	HTTPServer *http.Server
}

type HttpHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

func (c *Container) Init() error {
	if err := c.registerProducers(); err != nil {
		return err
	}

	if err := c.registerPaymentsConsumer(); err != nil {
		return err
	}

	if err := c.registerMoviesConsumer(); err != nil {
		return err
	}

	if err := c.registerUsersConsumer(); err != nil {
		return err
	}
	c.registerHttpHandlers()
	c.HTTPServer = &http.Server{
		Addr:    ":" + c.Config.Port,
		Handler: nil,
	}

	return nil
}

func (c *Container) registerHttpHandlers() {
	c.PaymentsHandler = server.NewPaymentsHandler(c.PaymentsProducer)
	c.MoviesHandler = server.NewMoviesHandler(c.MoviesProducer)
	c.UsersHandler = server.NewUsersHandler(c.UserProducer)
}

func (c *Container) registerProducers() error {
	var err error
	c.Producer, err = registerKafkaProducer(c.Config.ProducerConfig)
	if err != nil {
		return err
	}

	c.PaymentsProducer = kafkaAdapter.NewProducer(
		c.Config.PaymentsTopic,
		c.Producer, 
		func(p entities.Payment) []byte {
			 return []byte(strconv.Itoa(p.ID))
			},
		)

	c.MoviesProducer = kafkaAdapter.NewProducer(
		c.Config.MoviesTopic,
		c.Producer, 
		func(m entities.Movie) []byte {
				return []byte(strconv.Itoa(m.ID))
			},
		)

	c.UserProducer = kafkaAdapter.NewProducer(
		c.Config.UsersTopic,
		c.Producer, 
		func(u entities.User) []byte {
			 return []byte(strconv.Itoa(u.ID))
			},
		)

	return nil
}

func (c *Container) registerPaymentsConsumer() error {
	consumer, err := registerKafkaConsumer(c.Config.PaymentsConsumerConfig)
	if err != nil {
		return err
	}

	c.PaymentsConsumer = kafkaAdapter.NewConsumer(
		c.Config.PaymentsTopic, 
		c.Config.ConsumerPollTimeout,
		consumer,
		kafkaAdapter.NewKafkaProcessor(&services.PaymentsProcessor{}),
	)

	return nil
}

func (c *Container) registerMoviesConsumer() error {
	consumer, err := registerKafkaConsumer(c.Config.MoviesConsumerConfig)
	if err != nil {
		return err
	}

	c.MoviesConsumer = kafkaAdapter.NewConsumer(
		c.Config.MoviesTopic, 
		c.Config.ConsumerPollTimeout,
		consumer,
		kafkaAdapter.NewKafkaProcessor(&services.MoviesProcessor{}),
	)

	return nil
}

func (c *Container) registerUsersConsumer() error {
	consumer, err := registerKafkaConsumer(c.Config.UsersConsumerConfig)
	if err != nil {
		return err
	}

	c.UsersConsumer = kafkaAdapter.NewConsumer(
		c.Config.UsersTopic, 
		c.Config.ConsumerPollTimeout,
		consumer,
		kafkaAdapter.NewKafkaProcessor(&services.UsersProcessor{}),
	)

	return nil
}