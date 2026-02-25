package main

import (
	"fmt"
	"log"
	shared_kafka "micros/shared-kafka"
	"sync"
)

var addresses = []string{"localhost:9092", "localhost:9093", "localhost:9094"}
var topics = []string{"my-topic"}
var group = "my-group"

func main() {

	cons, err := shared_kafka.NewConsumer(addresses, topics, group)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	fmt.Println("starting consumer")
	go cons.StartConsuming()
	wg.Wait()
}
