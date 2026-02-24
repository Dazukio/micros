package models

import (
	"time"
)

type Task struct {
	Id       int       `json:"id"`
	Title    string    `json:"title"`
	Deadline time.Time `json:"deadline"` // будет как число
}
