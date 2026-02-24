package repository

import (
	"log"
	"micros/task-service/internal/models"
)

type Repository struct {
	taskList []*models.Task
}

func NewRepository() (*Repository, error) {
	return &Repository{taskList: make([]*models.Task, 0)}, nil
}

func (r *Repository) FindAll() []*models.Task {
	log.Print(r.taskList)

	return r.taskList
}

func (r *Repository) FindById(id int) *models.Task {
	for _, task := range r.taskList {
		if task.Id == id {
			return task
		}
	}
	return nil
}

func (r *Repository) Create(task *models.Task) {
	r.taskList = append(r.taskList, task)
}

func (r *Repository) Delete(id int) {
	r.taskList = append(r.taskList[:id], r.taskList[id+1:]...)
}
