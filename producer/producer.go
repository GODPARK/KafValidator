package producer

import (
	"errors"
	"fmt"

	"github.com/GODPARK/KafValidator/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProducerRunner struct {
	Runner   *kafka.Producer
	Topic    string
	MsgCount int
}

func KafkaProducerInit(config *config.Config) (*ProducerRunner, error) {
	producerRunner := &ProducerRunner{}
	runner, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.GetBootStrapServer(),
		"client.id":         config.Producer.ClientID,
		"acks":              config.Producer.Acks})
	if err != nil {
		return nil, errors.New("Kafka Producer Init Error: " + err.Error())
	}
	producerRunner.Runner = runner
	producerRunner.Topic = config.Topic
	producerRunner.MsgCount = config.Simple.MsgCount
	return producerRunner, nil
}

func (producer *ProducerRunner) Pub(value string) {
	delivery_chan := make(chan kafka.Event, 10000)
	for i := 0; i < producer.MsgCount; i++ {
		producer.Runner.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &producer.Topic, Partition: kafka.PartitionAny},
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
	}
	producer.Runner.Flush(1000)
	close(delivery_chan)
}
