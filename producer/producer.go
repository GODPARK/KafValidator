package producer

import (
	"errors"
	"fmt"

	"github.com/GODPARK/KafValidator/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProducerRunner struct {
	runner *kafka.Producer
	topic  string
}

func (producer *ProducerRunner) KafkaProducerInit(config *config.Config) error {
	producerRunner := &ProducerRunner{}
	fmt.Print(config.GetBootStrapServer())
	runner, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.GetBootStrapServer(),
		"client.id":         config.Producer.ClientID,
		"acks":              config.Producer.Acks})
	if err != nil {
		return errors.New("Kafka Producer Init Error: " + err.Error())
	}
	producerRunner.runner = runner
	producer.topic = config.Topic
	return nil
}

func (producer *ProducerRunner) Send(value string) {
	delivery_chan := make(chan kafka.Event, 10000)
	producer.runner.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &producer.topic, Partition: kafka.PartitionAny},
		Value:          []byte(value)},
		delivery_chan)

	e := <-delivery_chan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	close(delivery_chan)
}
