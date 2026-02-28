package shared_kafka

import (
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/OutBox/*.sql
var OutBoxMigrations embed.FS

type OutBox struct {
	DB        *sql.DB
	notifier  chan struct{}
	insertSql string
}

func NewOutBox(db *sql.DB) *OutBox {
	return &OutBox{DB: db, notifier: make(chan struct{}, 1), insertSql: "INSERT INTO messages(message, processed) VALUES (?,?)"}
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

func (ob *OutBox) WriteEvent(tx *sql.Tx, task *Task) error {
	select {
	case ob.notifier <- struct{}{}:
	default:
	}
	stmt, err := tx.Prepare(ob.insertSql)
	if err != nil {
		log.Println("Error preparing insert statement:", err)
		return err
	}
	defer stmt.Close()
	taskJson, err := json.Marshal(task)
	if err != nil {
		log.Println("Error marshalling task:", err)
		return err
	}
	_, err = stmt.Exec(string(taskJson), false)
	if err != nil {
		log.Println("Error inserting message into database:", err)
		return err
	}
	return nil
}

func (ob *OutBox) ReadEvents() (*Task, error) {
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
