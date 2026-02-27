package main

import (
	"context"
	"fmt"
	"log"
	shared_kafka "micros/shared-kafka"
	"micros/task-service/internal/handlers"
	"micros/task-service/internal/repository"
	"micros/task-service/internal/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type OutBoxReader interface {
	ReadEvents() (*shared_kafka.Task, error)
	GetNotifier() <-chan struct{}
}

var addresses = []string{"localhost:9092", "localhost:9093", "localhost:9094"}
var topic = "my-topic"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	repo := &repository.Repository{}
	OutboxWriter := shared_kafka.NewInMemory()
	StartWritingMessages(context.TODO(), OutboxWriter, addresses, topic)
	s := service.NewService(repo, OutboxWriter)
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

func processMessages(OB OutBoxReader, cons *shared_kafka.Producer) {
	for {
		msg, err := OB.ReadEvents()
		if err != nil {
			log.Println(err)
			break
		}
		if msg == nil {
			break
		}
		err = cons.Produce(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
