package service

import (
	"log"
	"micros/shared-kafka"
)

type OutboxWriter interface {
	WriteEvent(task *shared_kafka.Task) error
}

type Repository interface {
	FindById(id int) *shared_kafka.Task
	Create(task *shared_kafka.Task) error
	Delete(id int)
	FindAll() []*shared_kafka.Task
}

//type EventProducer interface {
//	Produce(task *shared_kafka.Task) error
//}

type Service struct {
	Repo Repository
	OW   OutboxWriter
}

func NewService(repo Repository, writer OutboxWriter) *Service {
	return &Service{Repo: repo, OW: writer}
}

func (s *Service) FindAll() []*shared_kafka.Task {
	return s.Repo.FindAll()
}

func (s *Service) FindById(id int) *shared_kafka.Task {
	return s.Repo.FindById(id)
}
func (s *Service) Create(task *shared_kafka.Task) error {
	s.Repo.Create(task)
	err := s.OW.WriteEvent(task)
	if err != nil {
		log.Println("Error producing task:", err)
	}
	return err
}
func (s *Service) Delete(id int) {
	s.Repo.Delete(id)
}
