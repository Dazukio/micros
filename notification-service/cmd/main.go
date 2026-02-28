package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"micros/notification-service/internal/repository"
	shared_kafka "micros/shared-kafka"
	"sync"
)

var addresses = []string{"localhost:9092", "localhost:9093", "localhost:9094"}
var topics = []string{"my-topic"}
var group = "my-group"
var DedupDBPath = "../internal/storage/db.db"

func main() {
	DedupDB, err := sql.Open("sqlite3", DedupDBPath)
	if err != nil {
		log.Fatal(err)
	}
	Consumer := repository.NewConsumer(DedupDB)
	cons, err := shared_kafka.NewConsumer(addresses, topics, group, Consumer)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	fmt.Println("starting consumer")
	go cons.StartConsuming()
	wg.Wait()
}
