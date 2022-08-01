package util

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/GODPARK/KafValidator/config"
)

type Context struct {
	BootstrapServer string
	Topic           string
	MsgCount        int
	Consumer        struct {
		ReciveTimeDiff []int64
		Average        int
		SuccessCount   int
		FailCount      int
		EtcCount       int
		MaxTime        int
		MinTime        int
	}
	Producer struct {
		SendTimeDiff []int64
		Average      int
		SuccessCount int
		FailCount    int
		EtcCount     int
		MaxTime      int
		MinTime      int
	}
}

func NewContext(config *config.Config) *Context {
	context := &Context{}

	context.BootstrapServer = config.GetBootStrapServer()
	context.Topic = config.Topic
	context.MsgCount = config.Simple.MsgCount

	context.Consumer.MinTime = math.MaxInt32
	context.Consumer.MaxTime = 0

	context.Producer.MinTime = math.MaxInt32
	context.Producer.MaxTime = 0

	return context
}

func (context *Context) AppendConsumerRDiff(timestamp int64) {
	now := (time.Now().UnixNano() / 1000000)
	diff := now - timestamp
	context.MaxTimeConsumer(int(diff))
	context.MinTimeConsumer(int(diff))
	context.Consumer.ReciveTimeDiff = append(context.Consumer.ReciveTimeDiff, now-timestamp)
}

func (context *Context) AppendProducerSDiff(timestamp int64) {
	now := (time.Now().UnixNano() / 1000000)
	diff := now - timestamp
	context.MaxTimeProducer(int(diff))
	context.MinTimeProducer(int(diff))
	context.Producer.SendTimeDiff = append(context.Producer.SendTimeDiff, now-timestamp)
}

func (context *Context) MaxTimeConsumer(value int) {
	if value >= context.Consumer.MaxTime {
		context.Consumer.MaxTime = value
	}
}

func (context *Context) MaxTimeProducer(value int) {
	if value >= context.Producer.MaxTime {
		context.Producer.MaxTime = value
	}
}

func (context *Context) MinTimeConsumer(value int) {
	if value < context.Consumer.MinTime {
		context.Consumer.MinTime = value
	}
}

func (context *Context) MinTimeProducer(value int) {
	if value < context.Producer.MinTime {
		context.Producer.MinTime = value
	}
}

func (context *Context) CalcConsuemrAvg() {
	var tmp int64 = 0
	for i := 0; i < len(context.Consumer.ReciveTimeDiff); i++ {
		tmp += context.Consumer.ReciveTimeDiff[0]
	}
	if tmp == 0 {
		context.Consumer.Average = -1
	} else {
		avg := int(tmp) / len(context.Consumer.ReciveTimeDiff)
		context.Consumer.Average = avg
	}
}

func (context *Context) CalcProducerAvg() {
	var tmp int64 = 0
	for i := 0; i < len(context.Producer.SendTimeDiff); i++ {
		tmp += context.Producer.SendTimeDiff[0]
	}
	if tmp == 0 {
		context.Producer.Average = -1
	} else {
		avg := int(tmp) / len(context.Producer.SendTimeDiff)
		context.Producer.Average = avg
	}
}

func setGreen(value string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", value)
}

func setRed(value string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", value)
}

func setCyan(value string) string {
	return fmt.Sprintf("\033[36m%s\033[0m", value)
}

func (context *Context) ShowResult() {
	context.CalcProducerAvg()
	context.CalcConsuemrAvg()

	fmt.Printf("\n\n")
	fmt.Printf("################### " + setCyan("Reslut") + " ####################\n")
	fmt.Printf("[CONFIG] Kafka Broker: %s\n", context.BootstrapServer)
	fmt.Printf("[CONFIG] Target Topic: %s\n", context.Topic)
	fmt.Printf("[CONFIG] Total Msg Count: %d\n", context.MsgCount)
	fmt.Print("\n")
	fmt.Printf("[PRODUCER] Msg %s Count: %s\n", setGreen("Pub Success"), setGreen(strconv.Itoa(context.Producer.SuccessCount)))
	fmt.Printf("[PRODUCER] Msg %s Count: %s\n", setRed("Pub Fail"), setRed(strconv.Itoa(context.Producer.FailCount)))
	fmt.Printf("[PRODUCER] Send Msg Average: %s ms\n", setGreen(strconv.Itoa(context.Producer.Average)))
	fmt.Printf("[PRODUCER] Send Msg Max Time: %d ms\n", context.Producer.MaxTime)
	fmt.Printf("[PRODUCER] Send Msg Min Time: %d ms\n", context.Producer.MinTime)
	fmt.Print("\n")
	fmt.Printf("[CONSUMER] Msg %s Count: %s\n", setGreen("Sub Success"), setGreen(strconv.Itoa(context.Consumer.SuccessCount)))
	fmt.Printf("[CONSUMER] Msg %s Count: %s\n", setRed("Sub Fail"), setRed(strconv.Itoa(context.Consumer.FailCount)))
	fmt.Printf("[CONSUMER] Msg %s Count: %s\n", setRed("Sub Etc"), setRed(strconv.Itoa(context.Consumer.EtcCount)))
	fmt.Printf("[CONSUMER] Pub->Sub Msg Time Average: %s ms\n", setGreen(strconv.Itoa(context.Consumer.Average)))
	fmt.Printf("[CONSUMER] Pub->Sub Msg Time Max: %d ms\n", context.Consumer.MaxTime)
	fmt.Printf("[CONSUMER] Pub->Sub Msg Time Min: %d ms\n", context.Consumer.MinTime)
	fmt.Printf("#################################################\n")
}
