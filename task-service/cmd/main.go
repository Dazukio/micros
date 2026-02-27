package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	shared_kafka "micros/shared-kafka"
	"micros/task-service/internal/handlers"
	"micros/task-service/internal/repository"
	"micros/task-service/internal/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

type OutBoxReader interface {
	ReadEvents() (*shared_kafka.Task, error)
	GetNotifier() <-chan struct{}
}

var addresses = []string{"localhost:9092", "localhost:9093", "localhost:9094"}
var topic = "my-topic"
var path = ".mydb.db"

func main() {
	r := chi.NewRouter()
	repo := &repository.Repository{}
	OutboxDB, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	defer OutboxDB.Close()
	err = shared_kafka.MigrateOutBox(OutboxDB)
	if err != nil {
		log.Fatal(err)
	}
	Outbox := shared_kafka.NewOutBox(OutboxDB)

	StartWritingMessages(context.TODO(), Outbox, addresses, topic)
	s := service.NewService(repo, Outbox)
	handler := handlers.NewHandler(s)
	r.Get("/tasks", handler.GetTasks)
	r.Post("/tasks", handler.AddTask)
	log.Println("Listening on port 8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}

}

func StartWritingMessages(ctx context.Context, OB OutBoxReader, addresses []string, topic string) {
	prod, err := shared_kafka.NewProducer(addresses, topic)
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(1 * time.Second)
	notifier := OB.GetNotifier()
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Closing writing messages")
				return
			case <-ticker.C:
				fmt.Println("tick check")
				processMessages(OB, prod)
			case <-notifier:
				fmt.Println("notifier check")
				processMessages(OB, prod)
			}
		}
	}()
}

func processMessages(OB OutBoxReader, prod *shared_kafka.Producer) {
	for {
		msg, err := OB.ReadEvents()
		if err != nil {
			log.Printf("Error reading events: %v", err)
			break
		}
		if msg == nil {
			break
		}

		log.Printf("Sending task %d to Kafka", msg.Id)
		if err := prod.Produce(msg); err != nil {
			log.Printf("Failed to send task %d: %v", msg.Id, err)
			continue
		}
		log.Printf("Task %d sent successfully", msg.Id)
	}
}
