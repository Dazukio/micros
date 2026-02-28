package repository

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	mu *sync.RWMutex
	db *sql.DB
}

func NewConsumer(db *sql.DB) *Consumer {
	return &Consumer{mu: &sync.RWMutex{}, db: db}
}

func (consumer *Consumer) HandleMessage(msg *kafka.Message) error {
	consumer.mu.Lock()
	defer consumer.mu.Unlock()
	getSql := consumer.db.QueryRow("SELECT * FROM messages WHERE ID = ?", msg.Key)
	err := getSql.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	return consumer.MarkProcessed(msg)
}

func (consumer *Consumer) MarkProcessed(msg *kafka.Message) error {
	id, err := strconv.Atoi(string(msg.Key))
	if err != nil {
		return err
	}
	_, err = consumer.db.Exec("INSERT INTO  messages (ID, message, processed) VALUES (?,?,?)", id, msg.Value, true)

	if err != nil {
		return err
	}
	log.Printf("%+v\n", string(msg.Value))
	return nil
}
