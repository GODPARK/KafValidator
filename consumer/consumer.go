package consumer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	Result   struct {
		diff    []int64
		average int
		success int
		fail    int
		etc     int
	}
}

func KafkaConsumerInit(config *config.Config, key *util.MsgKey) (*ConsumerRunner, error) {
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
						now := (time.Now().UnixNano() / 1000000)
						consumer.Result.diff = append(consumer.Result.diff, now-int64(timestamp))
						successMsgCount++
					}
					subMsgCount++
				} else {
					etcCount++
				}
				consumer.Runner.Commit()

				if subMsgCount >= consumer.MsgCount {
					consumer.Result.success = successMsgCount
					consumer.Result.fail = consumer.MsgCount - successMsgCount
					consumer.Result.etc = etcCount
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

func (consumer *ConsumerRunner) ShowResult() {

	var tmp int64 = 0
	for i := 0; i < len(consumer.Result.diff); i++ {
		tmp += consumer.Result.diff[0]
	}
	if tmp == 0 {
		consumer.Result.average = -1
	} else {
		avg := int(tmp) / len(consumer.Result.diff)
		consumer.Result.average = avg
	}

	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"

	fmt.Printf("[CONSUMER] Msg %sSub Success%s Count: %s%d%s\n",
		colorGreen, colorReset, colorGreen, consumer.Result.success, colorReset)
	fmt.Printf("[CONSUMER] Msg %sSub Fail%s Count: %s%d%s\n",
		colorRed, colorReset, colorRed, consumer.Result.fail, colorReset)
	fmt.Printf("[CONSUMER] Msg %sSub Other%s Count: %s%d%s\n",
		colorRed, colorReset, colorRed, consumer.Result.etc, colorReset)
	fmt.Printf("[CONSUMER] Msg Pub->Sub Time(average): %s%d%s ms\n", colorGreen, consumer.Result.average, colorReset)
}
