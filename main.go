package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/GODPARK/KafValidator/config"
	"github.com/GODPARK/KafValidator/consumer"
	"github.com/GODPARK/KafValidator/producer"
)

func main() {
	configPath := flag.String("config", "", "config.toml file path")
	flag.Parse()
	configData, err := config.InitConfig(*configPath)
	if err != nil {
		panic("Config Error Please check your config --> " + err.Error())
	}

	var wait sync.WaitGroup
	wait.Add(2)

	cInit, err := consumer.KafkaConsumerInit(configData)
	if err != nil {
		panic("Consumer Init Error: " + err.Error())
	}
	pInit, err := producer.KafkaProducerInit(configData)
	if err != nil {
		panic("Producer Init Error: " + err.Error())
	}

	go func(p *producer.ProducerRunner) {
		defer wait.Done()
		fmt.Println("producer start")
		p.Pub(time.Now().String())
	}(pInit)

	go func(c *consumer.ConsumerRunner) {
		defer wait.Done()
		fmt.Println("consumer start")
		c.Sub()
	}(cInit)

	wait.Wait()

}
