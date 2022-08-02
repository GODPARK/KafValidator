package main

/**
@author Park Jeong Hyun
@email rootjh937dev@gmail.com
*/

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/GODPARK/KafValidator/config"
	"github.com/GODPARK/KafValidator/consumer"
	"github.com/GODPARK/KafValidator/producer"
	"github.com/GODPARK/KafValidator/util"
)

func main() {
	configPath := flag.String("config", "", "config.toml file path")
	flag.Parse()
	configData, err := config.InitConfig(*configPath)
	if err != nil {
		panic("Config Error Please check your config --> " + err.Error())
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	TestKey := util.NewUUID()
	Context := util.NewContext(configData)
	cInit, err := consumer.KafkaConsumerInit(configData, TestKey, Context)
	if err != nil {
		panic("Consumer Init Error: " + err.Error())
	}
	pInit, err := producer.KafkaProducerInit(configData, TestKey, Context)
	if err != nil {
		panic("Producer Init Error: " + err.Error())
	}

	go func() {
		sig := <-sigs
		fmt.Printf("\nGET interrupt: %s -> producer , consumer close work... waiting please\n", sig.String())
		pInit.Runner.Flush(1000)
		cInit.Runner.Commit()
		pInit.Runner.Close()
		cInit.Runner.Close()
		fmt.Printf("close job is success.. exit!\n")
		Context.ShowResult()
		os.Exit(1)
	}()

	fmt.Printf("\033[32m Running...\033[0m please waiting\n")
	var wait sync.WaitGroup
	wait.Add(2)

	go func(p *producer.ProducerRunner) {
		defer wait.Done()
		p.Pub()
	}(pInit)

	go func(c *consumer.ConsumerRunner) {
		defer wait.Done()
		c.Sub()
	}(cInit)

	wait.Wait()

	pInit.Runner.Flush(1000)
	cInit.Runner.Commit()
	pInit.Runner.Close()
	cInit.Runner.Close()

	Context.ShowResult()

}
