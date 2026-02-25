package shared_kafka

import (
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	cons    *kafka.Consumer
	Running bool
}

func NewConsumer(addresses, topics []string, groupId string) (*Consumer, error) {
	conf := kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(addresses, ","),
		"group.id":           groupId,
		"auto.offset.reset":  "latest",
		"enable.auto.commit": false,
	}
	cons, err := kafka.NewConsumer(&conf)
	if err != nil {
		return nil, err
	}
	if err := cons.SubscribeTopics(topics, nil); err != nil {
		cons.Close()
		return nil, err
	}
	return &Consumer{cons: cons, Running: true}, nil
}

func (c *Consumer) StartConsuming() {
	for c.Running {
		msg, err := c.cons.ReadMessage(1000)
		if err != nil {
			kafkaErr, ok := err.(kafka.Error)
			if ok && kafkaErr.IsTimeout() {
				continue
			}
			log.Println("error reading from kafka:", err)
			continue
		}
		err = c.ProcessMessage(msg)
		if err != nil {

			log.Println("error processing message:", err)
			continue
		} else {
			c.cons.CommitMessage(msg)
		}
	}
}

func (c *Consumer) ProcessMessage(message *kafka.Message) error {
	log.Println(string(message.Value))
	return nil
}

func (c *Consumer) Stop() error {
	c.Running = false
	return c.cons.Close()
}
