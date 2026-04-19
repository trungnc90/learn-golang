package todo

import "time"

// Task represents a todo item.
type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
}

// storer defines the interface for task persistence.
// Any storage backend (file, database, memory) can implement this.
//go:generate go-mock-gen --interface=storer
type storer interface {
	Load() ([]Task, error)
	Save(tasks []Task) error
}
