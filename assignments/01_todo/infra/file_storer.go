package infra

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo"
)

// FileStore implements Storer using a JSON file.
type FileStore struct {
	FilePath string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{FilePath: path}
}

// load reads all tasks from the JSON file.
func (fs *FileStore) load() ([]todo.Task, error) {
	data, err := os.ReadFile(fs.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []todo.Task{}, nil
		}
		return nil, err
	}

	var tasks []todo.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// save writes all tasks to the JSON file.
func (fs *FileStore) save(tasks []todo.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fs.FilePath, data, 0644)
}

// nextID returns the next available ID based on existing tasks.
func nextID(tasks []todo.Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.Id > maxID {
			maxID = t.Id
		}
	}
	return maxID + 1
}

func (fs *FileStore) Create(task todo.Task) (todo.Task, error) {
	tasks, err := fs.load()
	if err != nil {
		return todo.Task{}, fmt.Errorf("load tasks: %w", err)
	}

	task.Id = nextID(tasks)
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}

	tasks = append(tasks, task)
	if err := fs.save(tasks); err != nil {
		return todo.Task{}, fmt.Errorf("save tasks: %w", err)
	}

	return task, nil
}

func (fs *FileStore) List(filter string) ([]todo.Task, error) {
	tasks, err := fs.load()
	if err != nil {
		return nil, fmt.Errorf("load tasks: %w", err)
	}

	if filter == "" {
		return tasks, nil
	}

	var filtered []todo.Task
	for _, t := range tasks {
		if filter == "done" && t.Done {
			filtered = append(filtered, t)
		}
		if filter == "pending" && !t.Done {
			filtered = append(filtered, t)
		}
	}
	if filtered == nil {
		filtered = []todo.Task{}
	}
	return filtered, nil
}

func (fs *FileStore) GetByID(id int) (todo.Task, error) {
	tasks, err := fs.load()
	if err != nil {
		return todo.Task{}, fmt.Errorf("load tasks: %w", err)
	}

	for _, t := range tasks {
		if t.Id == id {
			return t, nil
		}
	}
	return todo.Task{}, fmt.Errorf("task #%d not found", id)
}

func (fs *FileStore) Update(task todo.Task) (todo.Task, error) {
	tasks, err := fs.load()
	if err != nil {
		return todo.Task{}, fmt.Errorf("load tasks: %w", err)
	}

	for i, t := range tasks {
		if t.Id == task.Id {
			if task.Title != "" {
				tasks[i].Title = task.Title
			}
			if task.Description != "" {
				tasks[i].Description = task.Description
			}
			if task.Priority != "" {
				tasks[i].Priority = task.Priority
			}

			if err := fs.save(tasks); err != nil {
				return todo.Task{}, fmt.Errorf("save tasks: %w", err)
			}
			return tasks[i], nil
		}
	}
	return todo.Task{}, fmt.Errorf("task #%d not found", task.Id)
}

func (fs *FileStore) Delete(id int) error {
	tasks, err := fs.load()
	if err != nil {
		return fmt.Errorf("load tasks: %w", err)
	}

	for i, t := range tasks {
		if t.Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			if err := fs.save(tasks); err != nil {
				return fmt.Errorf("save tasks: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("task #%d not found", id)
}

func (fs *FileStore) ToggleDone(id int) (todo.Task, error) {
	tasks, err := fs.load()
	if err != nil {
		return todo.Task{}, fmt.Errorf("load tasks: %w", err)
	}

	for i, t := range tasks {
		if t.Id == id {
			tasks[i].Done = !tasks[i].Done
			if err := fs.save(tasks); err != nil {
				return todo.Task{}, fmt.Errorf("save tasks: %w", err)
			}
			return tasks[i], nil
		}
	}
	return todo.Task{}, fmt.Errorf("task #%d not found", id)
}
