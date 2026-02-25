package service

import (
	"log"
	models "micros/shared-kafka"
)

type Repository interface {
	FindById(id int) *models.Task
	Create(task *models.Task) error
	Delete(id int)
	FindAll() []*models.Task
}

type EventProducer interface {
	Produce(task *models.Task) error
}

type Service struct {
	Repo Repository
	EP   EventProducer
}

func NewService(repo Repository, EP EventProducer) *Service {
	return &Service{Repo: repo, EP: EP}
}

func (s *Service) FindAll() []*models.Task {
	return s.Repo.FindAll()
}

func (s *Service) FindById(id int) *models.Task {
	return s.Repo.FindById(id)
}
func (s *Service) Create(task *models.Task) error {
	s.Repo.Create(task)
	err := s.EP.Produce(task)
	if err != nil {
		log.Println("Error producing task:", err)
	}
	return err
}
func (s *Service) Delete(id int) {
	s.Repo.Delete(id)
}
