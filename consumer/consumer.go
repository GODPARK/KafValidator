package consumer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/GODPARK/KafValidator/config"
	"github.com/GODPARK/KafValidator/env"
	"github.com/GODPARK/KafValidator/util"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerRunner struct {
	Runner   *kafka.Consumer
	Topic    string
	MsgCount int
	Key      string
	Context  *util.Context
}

func KafkaConsumerInit(config *config.Config, key *util.MsgKey, context *util.Context) (*ConsumerRunner, error) {
	consumerRunner := &ConsumerRunner{}
	runner, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.GetBootStrapServer(),
		"group.id":          config.Consumer.GroupID,
		"auto.offset.reset": env.AUTO_OFFSET_RESET})
	if err != nil {
		return nil, errors.New("Kafka Consumer Init Error: " + err.Error())
	}
	consumerRunner.Runner = runner
	consumerRunner.Topic = config.Topic
	consumerRunner.MsgCount = config.Simple.MsgCount
	consumerRunner.Key = key.KeyToString()
	consumerRunner.Context = context
	return consumerRunner, nil
}

func (consumer *ConsumerRunner) Sub() {
	err := consumer.Runner.SubscribeTopics([]string{consumer.Topic}, nil)
	if err != nil {
		fmt.Println("Kafka Consumer Topic Subscribe Erro: " + err.Error())
	}
	run := true
	msgCount := 0
	subMsgCount := 0
	etcCount := 0
	successMsgCount := 0
	for run {
		ev := consumer.Runner.Poll(env.DEFAULT_POLL_INTERVAL_MS)
		switch e := ev.(type) {
		case *kafka.Message:
			if e != nil {
				msg := strings.Split(string(e.Value), ",")
				if msg[0] == consumer.Key {
					timestamp, err := strconv.Atoi(msg[1])
					if err == nil {
						consumer.Context.AppendConsumerRDiff(int64(timestamp))
						successMsgCount++
					}
					subMsgCount++
				} else {
					etcCount++
				}
				consumer.Runner.Commit()

				if subMsgCount >= consumer.MsgCount {
					consumer.Context.Consumer.SuccessCount = successMsgCount
					consumer.Context.Consumer.FailCount = consumer.MsgCount - successMsgCount
					consumer.Context.Consumer.EtcCount = etcCount
					run = false
				}
			}
		default:
			msgCount++
			if msgCount >= env.THRESHOLD_MSG_COUNT {
				run = false
			}
		}
	}
}
