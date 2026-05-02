package todo

import (
	"fmt"
	"time"
)

// MemoryStore implements Storer using an in-memory slice.
// Used for testing — no file I/O, no cleanup needed.
type MemoryStore struct {
	tasks  []Task
	nextID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{tasks: []Task{}, nextID: 1}
}

func (ms *MemoryStore) Create(task Task) (Task, error) {
	task.Id = ms.nextID
	ms.nextID++
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	ms.tasks = append(ms.tasks, task)
	return task, nil
}

func (ms *MemoryStore) List(filter string) ([]Task, error) {
	var result []Task
	for _, t := range ms.tasks {
		if filter == "done" && !t.Done {
			continue
		}
		if filter == "pending" && t.Done {
			continue
		}
		result = append(result, t)
	}
	if result == nil {
		result = []Task{}
	}
	return result, nil
}

func (ms *MemoryStore) GetByID(id int) (Task, error) {
	for _, t := range ms.tasks {
		if t.Id == id {
			return t, nil
		}
	}
	return Task{}, fmt.Errorf("task #%d not found", id)
}

func (ms *MemoryStore) Update(task Task) (Task, error) {
	for i, t := range ms.tasks {
		if t.Id == task.Id {
			if task.Title != "" {
				ms.tasks[i].Title = task.Title
			}
			if task.Description != "" {
				ms.tasks[i].Description = task.Description
			}
			if task.Priority != "" {
				ms.tasks[i].Priority = task.Priority
			}
			return ms.tasks[i], nil
		}
	}
	return Task{}, fmt.Errorf("task #%d not found", task.Id)
}

func (ms *MemoryStore) Delete(id int) error {
	for i, t := range ms.tasks {
		if t.Id == id {
			ms.tasks = append(ms.tasks[:i], ms.tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("task #%d not found", id)
}

func (ms *MemoryStore) ToggleDone(id int) (Task, error) {
	for i, t := range ms.tasks {
		if t.Id == id {
			ms.tasks[i].Done = !ms.tasks[i].Done
			return ms.tasks[i], nil
		}
	}
	return Task{}, fmt.Errorf("task #%d not found", id)
}
