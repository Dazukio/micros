package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"micros/task-service/internal/models"
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
	task := &models.Task{}
	err = json.Unmarshal(body, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	h.Service.Create(task)
	_, err = http.Post("http://127.0.0.1:8081/notify", "application/json", bytes.NewReader(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
