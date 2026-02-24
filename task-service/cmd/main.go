package main

import (
	"log"
	"micros/task-service/internal/handlers"
	"micros/task-service/internal/repository"
	"micros/task-service/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	repo := &repository.Repository{}
	s := service.NewService(repo)
	handler := handlers.NewHandler(s)
	r.Get("/tasks", handler.GetTasks)
	r.Post("/tasks", handler.AddTask)
	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
