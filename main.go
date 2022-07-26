package main

import (
	"flag"
	"fmt"
)

func main() {
	configPath := flag.String("config", "", "config.toml file path")
	flag.Parse()

	fmt.Println(*configPath)
}