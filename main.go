package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/GODPARK/KafValidator/config"
)

func main() {
	configPath := flag.String("config", "", "config.toml file path")
	flag.Parse()
	configObject := &config.Config{}
	configData, err := configObject.InitConfig(*configPath)
	if err != nil {
		fmt.Printf("Config Error Please check your config --> %s", err)
	}
	fmt.Println(configData)

	var wait sync.WaitGroup
	wait.Add(2)

	go func() {
		defer wait.Done()

	}()

	go func() {
		defer wait.Done()
	}()

	wait.Wait()

}
