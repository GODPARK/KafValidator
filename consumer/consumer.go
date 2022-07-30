package consumer

import (
	"errors"
	"fmt"

	"github.com/GODPARK/KafValidator/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerRunner struct {
	Runner *kafka.Consumer
	Topic  string
}

func KafkaConsumerInit(config *config.Config) (*ConsumerRunner, error) {
	consumerRunner := &ConsumerRunner{}
	runner, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.GetBootStrapServer(),
		"group.id":          config.Consumer.GroupID,
		"auto.offset.reset": config.Consumer.AutoOffsetReset})
	if err != nil {
		return nil, errors.New("Kafka Consumer Init Error: " + err.Error())
	}
	consumerRunner.Runner = runner
	consumerRunner.Topic = config.Topic
	return consumerRunner, nil
}

func (consumer *ConsumerRunner) Sub() {
	err := consumer.Runner.Subscribe(consumer.Topic, nil)
	if err != nil {
		fmt.Println("Kafka Consumer Topic Subscribe Erro: " + err.Error())
	}
	msg, err := consumer.Runner.ReadMessage(0)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(msg)
	consumer.Runner.Close()
}
