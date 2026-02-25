package main

import (
	"log"
	shared_kafka "micros/shared-kafka"
	"micros/task-service/internal/handlers"
	"micros/task-service/internal/repository"
	"micros/task-service/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var addresses = []string{"localhost:9092", "localhost:9093", "localhost:9094"}
var topic = "my-topic"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	repo := &repository.Repository{}
	EP := shared_kafka.NewProducer(addresses, topic)
	s := service.NewService(repo, EP)
	handler := handlers.NewHandler(s)
	r.Get("/tasks", handler.GetTasks)
	r.Post("/tasks", handler.AddTask)
	log.Println("Listening on port 8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}

}
