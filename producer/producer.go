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
	Result   struct {
		success int
		fail    int
	}
}

func KafkaProducerInit(config *config.Config, key *util.MsgKey) (*ProducerRunner, error) {
	producerRunner := &ProducerRunner{}
	runner, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.GetBootStrapServer(),
		"client.id":         config.Producer.ClientID,
		"acks":              config.Producer.Acks})
	if err != nil {
		fmt.Println("1" + err.Error())
		return nil, errors.New("Kafka Producer Init Error: " + err.Error())
	}
	producerRunner.Runner = runner
	producerRunner.Topic = config.Topic
	producerRunner.MsgCount = config.Simple.MsgCount
	producerRunner.Key = key.KeyToString()
	return producerRunner, nil
}

func (producer *ProducerRunner) Pub() {
	delivery_chan := make(chan kafka.Event, 10000)
	successMsgCount := 0
	for i := 0; i < producer.MsgCount; i++ {
		now := fmt.Sprint(time.Now().UnixNano() / 1000000)
		msg := producer.Key + "," + now
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
		}
	}
	producer.Result.success = successMsgCount
	producer.Result.fail = producer.MsgCount - successMsgCount
	producer.Runner.Flush(10)
	close(delivery_chan)
}

func (producer *ProducerRunner) ShowResult() {
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"

	fmt.Printf("[PRODUCER] Msg %sPub Success%s Count: %s%d%s\n",
		colorGreen, colorReset, colorGreen, producer.Result.success, colorReset)
	fmt.Printf("[PRODUCER] Msg %sPub Fail%s Count: %s%d%s\n",
		colorRed, colorReset, colorRed, producer.Result.fail, colorReset)
}
