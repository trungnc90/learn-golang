package main

// Store defines the interface for task persistence.
// Any storage backend (file, database, memory) can implement this.
type Store interface {
	Load() ([]Task, error)
	Save(tasks []Task) error
}
