package handlers

import (
	"encoding/json"
	"io"
	shared_kafka "micros/shared-kafka"
	"micros/task-service/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()
	task := &shared_kafka.Task{}
	err = json.Unmarshal(body, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = h.Service.Create(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *Handler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	task := h.Service.FindById(id)
	if task == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(task)
}
