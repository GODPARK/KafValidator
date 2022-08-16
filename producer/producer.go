package producer

import (
	"errors"
	"fmt"
	"time"

	"github.com/GODPARK/KafValidator/config"
	"github.com/GODPARK/KafValidator/util"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProducerRunner struct {
	Runner   *kafka.Producer
	Topic    string
	MsgCount int
	Key      string
	Dummy    string
	Context  *util.Context
}

func KafkaProducerInit(config *config.Config, key *util.MsgKey, context *util.Context, dummy string) (*ProducerRunner, error) {
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
	producerRunner.Key = key.KeyToString()
	producerRunner.Context = context
	producerRunner.Dummy = dummy
	return producerRunner, nil
}

func (producer *ProducerRunner) Pub() {
	delivery_chan := make(chan kafka.Event, 10000)
	successMsgCount := 0
	for i := 0; i < producer.MsgCount; i++ {
		now := (time.Now().UnixNano() / 1000000)
		nowStr := fmt.Sprint(now)
		msg := producer.Key + "," + nowStr + "," + producer.Dummy

		producer.Runner.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &producer.Topic, Partition: kafka.PartitionAny},
			Value:          []byte(msg)},
			delivery_chan)

		e := <-delivery_chan
		m := e.(*kafka.Message)

		if m.TopicPartition.Error != nil {
			fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		} else {
			successMsgCount++
			producer.Context.AppendProducerSDiff(now)
		}
	}
	producer.Context.Producer.SuccessCount = successMsgCount
	producer.Context.Producer.FailCount = producer.MsgCount - successMsgCount
	producer.Runner.Flush(10)
	close(delivery_chan)
}
