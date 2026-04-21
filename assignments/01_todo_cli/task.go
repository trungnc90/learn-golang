package todo

import (
	"fmt"
	"time"
)

func nextId(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.Id > maxID {
			maxID = t.Id
		}
	}
	return maxID + 1
}

func (t *Todo) AddTask(cmd *AddCmd) (*Task, error) {
	tasks, err := t.store.Load()
	if err != nil {
		return nil, fmt.Errorf("load tasks: %w", err)
	}

	priority := cmd.Priority
	if priority == "" {
		priority = "low"
	}

	task := Task{
		Id:          nextId(tasks),
		Title:       cmd.Title,
		Description: cmd.Description,
		Priority:    priority,
		Done:        false,
		CreatedAt:   time.Now(),
	}

	tasks = append(tasks, task)
	if err := t.store.Save(tasks); err != nil {
		return nil, fmt.Errorf("save tasks: %w", err)
	}

	return &task, nil
}

func (t *Todo) ListTasks(cmd *ListCmd) ([]Task, error) {
	tasks, err := t.store.Load()
	if err != nil {
		return nil, fmt.Errorf("load tasks: %w", err)
	}

	if cmd.Filter == "" {
		return tasks, nil
	}

	var filtered []Task
	for _, task := range tasks {
		if cmd.Filter == "done" && task.Done {
			filtered = append(filtered, task)
		}
		if cmd.Filter == "pending" && !task.Done {
			filtered = append(filtered, task)
		}
	}
	return filtered, nil
}

func (t *Todo) UpdateTasks(cmd *UpdateCmd) (*Task, error) {
	tasks, err := t.store.Load()
	if err != nil {
		return nil, fmt.Errorf("load tasks: %w", err)
	}

	for i, task := range tasks {
		if task.Id == cmd.Id {
			if cmd.Title != "" {
				tasks[i].Title = cmd.Title
			}
			if cmd.Description != "" {
				tasks[i].Description = cmd.Description
			}
			if cmd.Priority != "" {
				tasks[i].Priority = cmd.Priority
			}

			if err := t.store.Save(tasks); err != nil {
				return nil, fmt.Errorf("save tasks: %w", err)
			}
			return &tasks[i], nil
		}
	}
	return nil, fmt.Errorf("task #%d not found", cmd.Id)
}

func (t *Todo) DeleteTask(cmd *DeleteCmd) error {
	tasks, err := t.store.Load()
	if err != nil {
		return fmt.Errorf("load tasks: %w", err)
	}

	for i, task := range tasks {
		if task.Id == cmd.Id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			if err := t.store.Save(tasks); err != nil {
				return fmt.Errorf("save tasks: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("task #%d not found", cmd.Id)
}

func (t *Todo) ToggleDone(cmd *DoneCmd) (*Task, error) {
	tasks, err := t.store.Load()
	if err != nil {
		return nil, fmt.Errorf("load tasks: %w", err)
	}

	for i, task := range tasks {
		if task.Id == cmd.Id {
			tasks[i].Done = !tasks[i].Done
			if err := t.store.Save(tasks); err != nil {
				return nil, fmt.Errorf("save tasks: %w", err)
			}
			return &tasks[i], nil
		}
	}
	return nil, fmt.Errorf("task #%d not found", cmd.Id)
}
