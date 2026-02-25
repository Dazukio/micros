package repository

import (
	"log"
	shared_kafka "micros/shared-kafka"
)

type Repository struct {
	taskList []*shared_kafka.Task
}

func NewRepository() (*Repository, error) {
	return &Repository{taskList: make([]*shared_kafka.Task, 0)}, nil
}

func (r *Repository) FindAll() []*shared_kafka.Task {
	log.Print(r.taskList)

	return r.taskList
}

func (r *Repository) FindById(id int) *shared_kafka.Task {
	for _, task := range r.taskList {
		if task.Id == id {
			return task
		}
	}
	return nil
}

func (r *Repository) Create(task *shared_kafka.Task) error {
	r.taskList = append(r.taskList, task)
	return nil
}

func (r *Repository) Delete(id int) {
	r.taskList = append(r.taskList[:id], r.taskList[id+1:]...)
}
