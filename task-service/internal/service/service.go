package service

import (
	"database/sql"
	"log"
	"micros/shared-kafka"
)

type OutboxWriter interface {
	WriteEvent(tx *sql.Tx, task *shared_kafka.Task) error
}

type Repository interface {
	FindById(id int) *shared_kafka.Task
	Create(tx *sql.Tx, task *shared_kafka.Task) error
}

//type EventProducer interface {
//	Produce(task *shared_kafka.Task) error
//}

type Service struct {
	Repo Repository
	OW   OutboxWriter
	db   *sql.DB
}

func NewService(repo Repository, db *sql.DB, writer OutboxWriter) *Service {
	return &Service{Repo: repo, OW: writer, db: db}
}

func (s *Service) FindById(id int) *shared_kafka.Task {
	return s.Repo.FindById(id)
}
func (s *Service) Create(task *shared_kafka.Task) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}
	defer tx.Rollback()
	err = s.Repo.Create(tx, task)
	if err != nil {
		log.Println("Error creating task:", err)
		return err
	}
	err = s.OW.WriteEvent(tx, task)
	if err != nil {
		return err
	}
	return tx.Commit()
}
