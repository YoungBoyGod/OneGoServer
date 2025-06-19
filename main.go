package main

import (
	"fmt"
	"log"
)

func init() {
	config, err := ReadConfig("config/server.yaml")
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	fmt.Println(config)
}

func main() {

	fmt.Println("Hello, World!")
}
