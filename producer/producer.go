package producer

import (
	"github.com/GODPARK/KafValidator/config"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProducerRunner struct {
	runner interface{}
}

func (producer *ProducerRunner) KafkaProducerInit(config *Config) (*ProducerRunner) {
	producerRunner := &ProducerRunner{}
	producerRunner.runner, err := kafka.NewProducer(

	)

	return producerRunner
}
