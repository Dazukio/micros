package service

import (
	"micros/task-service/internal/models"
)

type Repository interface {
	FindById(id int) *models.Task
	Create(task *models.Task)
	Delete(id int)
	FindAll() []*models.Task
}

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) FindAll() []*models.Task {
	return s.Repo.FindAll()
}

func (s *Service) FindById(id int) *models.Task {
	return s.Repo.FindById(id)
}
func (s *Service) Create(task *models.Task) {
	s.Repo.Create(task)
}
func (s *Service) Delete(id int) {
	s.Repo.Delete(id)
}
