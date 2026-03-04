package app


type Config struct {
	PaymentsTopic string
	UsersTopic    string
	MoviesTopic   string

	Port string

	ConsumerPollTimeout int

	PaymentsConsumerConfig map[string]interface{}
	UsersConsumerConfig map[string]interface{}
	MoviesConsumerConfig map[string]interface{}
	ProducerConfig map[string]interface{}
}