package shared_kafka

import (
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/OutBox/*.sql
var OutBoxMigrations embed.FS

type OutBox struct {
	DB        *sql.DB
	mu        *sync.RWMutex
	notifier  chan struct{}
	insertSql string
}

func NewOutBox(db *sql.DB) *OutBox {
	return &OutBox{DB: db, notifier: make(chan struct{}, 1), mu: &sync.RWMutex{}, insertSql: "INSERT INTO messages(id, message, processed) VALUES (?,?,?)"}
}

func MigrateOutBox(db *sql.DB) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	goose.SetBaseFS(OutBoxMigrations)

	if err := goose.Up(db, "migrations/OutBox"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (ob *OutBox) WriteEvent(task *Task) error {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	select {
	case ob.notifier <- struct{}{}:
	default:
	}
	stmt, err := ob.DB.Prepare(ob.insertSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	taskJson, err := json.Marshal(task)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(task.Id, string(taskJson), false)
	if err != nil {
		log.Println("Error inserting message into database:", err)
		return err
	}
	return nil
}

func (ob *OutBox) ReadEvents() (*Task, error) {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	var taskId int
	var messageJson string
	err := ob.DB.QueryRow("SELECT id, message FROM messages WHERE processed = false LIMIT 1").Scan(&taskId, &messageJson)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var task Task
	if err = json.Unmarshal([]byte(messageJson), &task); err != nil {
		return nil, err
	}

	err = ob.markProcessed(taskId)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (ob *OutBox) GetNotifier() <-chan struct{} {
	return ob.notifier
}

func (ob *OutBox) markProcessed(id int) error {
	_, err := ob.DB.Exec("UPDATE messages SET processed = true WHERE id = ?", id)
	return err
}
