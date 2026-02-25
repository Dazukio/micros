package handlers

import (
	"encoding/json"
	"io"
	shared_kafka "micros/shared-kafka"
	"micros/task-service/internal/service"
	"net/http"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.Service.FindAll())
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()
	task := &shared_kafka.Task{}
	err = json.Unmarshal(body, task)
	task.Id = len(h.Service.FindAll())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = h.Service.Create(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
