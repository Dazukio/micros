package repository

import (
	"database/sql"
	shared_kafka "micros/shared-kafka"
	"time"
)

type Repository struct {
	db        *sql.DB
	insertSql string
	findBySql string
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db, insertSql: "INSERT INTO tasks (TITLE, DEADLINE) VALUES (?, ?);",
		findBySql: "SELECT ID, TITLE, DEADLINE  FROM tasks WHERE ID = ? LIMIT 1;"}
}

func (r *Repository) FindById(id int) *shared_kafka.Task {
	var taskId int
	var title string
	var deadline time.Time
	err := r.db.QueryRow(r.findBySql, id).Scan(&taskId, &title, &deadline)
	if err != nil {
		return nil
	}
	return &shared_kafka.Task{Id: id, Title: title, Deadline: deadline}
}

func (r *Repository) Create(tx *sql.Tx, task *shared_kafka.Task) error {
	result, err := tx.Exec(r.insertSql, task.Title, task.Deadline)
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.Id = int(id)
	return err
}
