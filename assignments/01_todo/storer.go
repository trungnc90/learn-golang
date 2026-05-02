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
	Create(task Task) (Task, error)
	List(filter string) ([]Task, error)
	GetByID(id int) (Task, error)
	Update(task Task) (Task, error)
	Delete(id int) error
	ToggleDone(id int) (Task, error)
}
