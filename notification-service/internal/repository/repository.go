package repository

import (
	"encoding/json"
	"fmt"
	"log"
	shared_kafka "micros/shared-kafka"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	mu              *sync.RWMutex
	processedEvents map[int]bool
}

func NewConsumer() *Consumer {
	return &Consumer{mu: &sync.RWMutex{}, processedEvents: make(map[int]bool)}
}

func (consumer *Consumer) HandleMessage(msg *kafka.Message) error {
	var task shared_kafka.Task
	err := json.Unmarshal(msg.Value, &task)
	if err != nil {
		return err
	}

	consumer.mu.RLock()
	if consumer.processedEvents[task.Id] {
		log.Println("Duplicated event", task.Id)
		consumer.mu.RUnlock()
		return nil
	}

	consumer.mu.RUnlock()

	fmt.Printf("New event: %#+v\n", task)

	consumer.mu.Lock()
	consumer.processedEvents[task.Id] = true
	consumer.mu.Unlock()
	return nil
}
