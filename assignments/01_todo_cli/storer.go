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

// Storer defines the interface for task persistence.
//go:generate go-mock-gen --interface=Storer
type Storer interface {
	Load() ([]Task, error)
	Save(tasks []Task) error
}
