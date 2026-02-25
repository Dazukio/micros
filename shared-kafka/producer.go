package shared_kafka

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	Prod         *kafka.Producer
	DeliveryChan chan kafka.Event
	Topic        string
	Running      bool
}

func NewProducer(addresses []string, topic string) (*Producer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers": strings.Join(addresses, ","),
	}
	prod, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, err
	}
	return &Producer{Prod: prod, Running: true, DeliveryChan: make(chan kafka.Event), Topic: topic}, nil
}

func (p *Producer) Produce(message *Task) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: bytes,
		Key:   []byte(strconv.Itoa(message.Id)),
	}
	err = p.Prod.Produce(msg, p.DeliveryChan)
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) Stop() {
	p.Running = false
	p.Prod.Flush(5000)
	close(p.DeliveryChan)
	p.Prod.Close()
}

func (p *Producer) CheckResponse() {
	go func() {
		for p.Running {
			select {
			case e := <-p.DeliveryChan:
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						log.Printf("Delivery failed: %v", ev.TopicPartition)
					} else {
						log.Printf("Success delivery message: %v", ev)
					}
				case kafka.Error:
					log.Printf("Delivery failed: %v", ev.Error())

				}
			}
		}
	}()
}
